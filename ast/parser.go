package ast

import (
	"os"
	"strings"
	"zyro/code"
	"zyro/code/mod"
	"zyro/code/wrapper"
	"zyro/logger"
	"zyro/token"
	"zyro/tokenUtil/splitter"
)

// The standard library's path
var stdPath = os.Getenv("ZYRO_STANDARD_PATH")

// pushParameter pushes the current parameter to the parameters list
func pushParameter(
	unit token.Token,
	currentParameter *[]token.Token,
	currentParameterName *string,
	currentFunctionParameters *map[string][]token.Token,
	expectingArgType *bool,
	expectingArgs *bool,
) {
	// No parameter name = No parameters, skip
	if len(*currentParameterName) == 0 {
		return
	}

	// Parse next parameter
	if len(*currentParameter) == 0 {
		logger.TokenError(
			unit,
			"Expected a data type",
			"Add a data type",
		)
	}

	(*currentFunctionParameters)[*currentParameterName] = *currentParameter
	*currentParameterName = ""
	*currentParameter = nil
	*currentParameter = []token.Token{}
	*expectingArgType = false
	*expectingArgs = true
}

// Parse parses the given tokens into a FileCode
func Parse(tokens []token.Token, allowMods bool, allowInlineVars bool) *FileCode {
	result := NewFileCode()

	// Used to keep track of the state of the parser
	inFunction := false
	inMod := false

	// Start at block depth 1 because of the curly brace
	// used to define the function's body
	blockDepth := 1

	// Used to keep track of the expected tokens
	expectingFun := true
	expectingFunName := false
	expectingOpenParen := false
	expectingArgs := false
	expectingArgColon := false
	expectingArgType := false
	expectingArrow := false
	expectingReturnType := false
	expectingOpenCurly := false

	// Used to keep track of the current function's metadata
	currentFunctionName := ""
	currentParameterName := ""
	var currentFunctionReturnType []token.Token
	currentFunctionParameters := make(map[string][]token.Token)
	var currentFunctionTrace token.Token

	var currentFunctionPublic bool
	currentParameter := make([]token.Token, 0)
	currentFunctionBody := make([]token.Token, 0)

	currentModVars := make([][]token.Token, 0)

	// Used to skip tokens
	skipToIndex := 0

	// Used to parse generics in mods
	inGeneric := false
	genericDepth := 0
	genericTokens := make([]token.Token, 0)

	for i, unit := range tokens {
		if i < skipToIndex {
			continue
		}

		tokenType := unit.GetType()

		if (tokenType == token.Let || tokenType == token.Const) && !inFunction {
			if !allowInlineVars && !inMod {
				logger.TokenError(
					unit,
					"Inline variable declarations are not allowed here",
					"Remove the inline variable declaration",
				)
			}

			// Extract the statement
			statement := splitter.ExtractTokensBefore(
				tokens[i:],
				token.Semicolon,
				// Don't handle nested statements here
				false,
				token.Unknown,
				token.Unknown,
				true,
			)

			currentModVars = append(currentModVars, statement)

			skipToIndex = i + len(statement) + 1
			continue
		}

		if tokenType == token.Pub {
			if inFunction {
				logger.TokenError(
					unit,
					"Cannot declare a public function inside another function",
					"Move the function outside of the current function",
				)
			}

			// Make sure there is a next token and that it's in the same file
			if i == len(tokens)-1 || tokens[i+1].GetFile() != unit.GetFile() {
				logger.TokenError(
					unit,
					"Expected a function declaration after the public keyword",
					"Add a function declaration after the public keyword",
				)
			}

			if currentFunctionPublic {
				logger.TokenError(
					unit,
					"Cannot declare a function as public twice",
					"Remove the extra public keyword",
				)
			}

			currentFunctionPublic = true
		} else if expectingFun {
			if tokenType == token.Mod {
				if !allowMods {
					logger.TokenError(
						unit,
						"Modules are not allowed here",
						"Remove the module declaration",
					)
				}

				inMod = true
				expectingFun = false
				expectingFunName = true
				continue
			}

			if tokenType != token.Function {
				logger.TokenError(
					unit,
					"Expected a function declaration",
					"Add a function declaration",
				)
			}

			inFunction = true
			expectingFun = false
			expectingFunName = true
		} else if expectingFunName {
			// Only allow identifiers as function names
			if tokenType != token.Identifier {
				logger.TokenError(
					unit,
					"Expected a function name",
					"Add a function name",
				)
			}

			currentFunctionTrace = unit
			currentFunctionName = unit.GetValue()

			if inMod {
				expectingOpenCurly = true
				expectingFunName = false
				continue
			}

			expectingFunName = false
			expectingOpenParen = true
		} else if expectingOpenCurly {
			// This case only happens for modules
			if tokenType == token.LessThan && !inGeneric && genericDepth == 0 {
				inGeneric = true
				genericDepth++
				continue
			} else if tokenType == token.GreaterThan && inGeneric {
				genericDepth--

				if genericDepth < 0 {
					logger.TokenError(
						unit,
						"Unexpected closing angle bracket",
						"Remove the closing angle bracket",
					)
				}

				if genericDepth == 0 {
					if len(genericTokens) == 0 {
						logger.TokenError(
							unit,
							"Expected at least one generic parameter",
							"Add at least one generic parameter",
						)
					}

					inGeneric = false
					continue
				}

				genericTokens = append(genericTokens, unit)
			} else if inGeneric {
				genericTokens = append(genericTokens, unit)
				continue
			}

			if tokenType != token.OpenCurly {
				logger.TokenError(
					unit,
					"Expected an opening curly brace",
					"Add an opening curly brace",
				)
			}

			expectingOpenCurly = false
			// Go directly to parse the body of the mod
		} else if expectingOpenParen {
			if tokenType != token.OpenParen {
				logger.TokenError(
					unit,
					"Expected an open parenthesis",
					"Add an open parenthesis",
				)
			}

			expectingOpenParen = false
			expectingArgs = true
		} else if expectingArgs {
			if tokenType == token.CloseParen {
				pushParameter(
					unit,
					&currentParameter,
					&currentParameterName,
					&currentFunctionParameters,
					&expectingArgType,
					&expectingArgs,
				)

				expectingArgs = false
				expectingArrow = true
				continue
			}

			// Take the parameter's name
			// only identifiers are allowed here
			if tokenType != token.Identifier {
				logger.TokenError(
					unit,
					"Expected a parameter name",
					"Add a parameter name",
				)
			}

			currentParameterName = unit.GetValue()
			expectingArgs = false
			expectingArgColon = true
		} else if expectingArgColon {
			if tokenType != token.Colon {
				logger.TokenError(
					unit,
					"Expected a colon after the parameter name",
					"Add a colon after the parameter name",
				)
			}

			expectingArgColon = false
			expectingArgType = true
		} else if expectingArgType {
			if tokenType == token.Comma {
				pushParameter(
					unit,
					&currentParameter,
					&currentParameterName,
					&currentFunctionParameters,
					&expectingArgType,
					&expectingArgs,
				)
				continue
			}

			if tokenType == token.CloseParen {
				// Parse the return type
				pushParameter(
					unit,
					&currentParameter,
					&currentParameterName,
					&currentFunctionParameters,
					&expectingArgType,
					&expectingArgs,
				)

				expectingArgs = false
				expectingArrow = true
				expectingArgType = false
				continue
			}

			currentParameter = append(currentParameter, unit)
		} else if expectingArrow {
			if tokenType == token.OpenCurly {
				returnType := token.NewToken(
					token.Nothing,
					"nothing",
					unit.GetFile(),
					unit.GetLine(),
					unit.GetColumn(),
					"nothing",
					0,
					strings.Builder{},
				)

				// Default to nothing if no return type is specified
				currentFunctionReturnType = append(
					currentFunctionReturnType,
					*returnType,
				)

				expectingArrow = false
				expectingReturnType = false
				continue
			}

			if tokenType != token.Arrow {
				logger.TokenError(
					unit,
					"Expected an arrow after the function parameters",
					"Add an arrow after the function parameters",
				)
			}

			expectingArrow = false
			expectingReturnType = true
		} else if expectingReturnType {
			if tokenType == token.OpenCurly {
				// End of return type
				expectingReturnType = false
				continue
			}

			currentFunctionReturnType = append(currentFunctionReturnType, unit)
		} else {
			if !inFunction && !inMod {
				logger.TokenError(
					unit,
					"Unexpected token",
					"Remove the token",
					"You can't have a token outside of a function",
				)
			}

			if tokenType == token.OpenCurly {
				// Start of function body
				blockDepth++
			} else if tokenType == token.CloseCurly {
				// End of function body
				blockDepth--

				if blockDepth < 0 {
					logger.TokenError(
						unit,
						"Unexpected closing curly brace",
						"Remove the closing curly brace",
					)
				}

				if blockDepth == 0 {
					// End of function
					blockDepth = 1
					inFunction = false

					if inMod {
						// Recursively parse the mod
						// No risk of exponential complexity because
						// we don't allow nested mods
						modFunctions := Parse(currentFunctionBody, false, true)
						privateFunctions := make(map[string]*code.Function)
						publicFunctions := make(map[string]*code.Function)

						for _, functions := range *modFunctions.GetFunctions() {
							for name, function := range functions {
								if function.IsPublic() {
									publicFunctions[name] = function
								} else {
									privateFunctions[name] = function
								}
							}
						}

						// Split and parse generics
						genericParamsTokens := splitter.SplitTokens(
							genericTokens,
							token.Comma,
							token.LessThan,
							token.GreaterThan,
						)

						genericParams := make([]wrapper.TypeWrapper, len(genericParamsTokens))
						for n, paramTokens := range genericParamsTokens {
							genericParams[n] = wrapper.NewTypeWrapper(
								paramTokens,
								unit,
							)
						}

						// Wrap the functions inside a ZyroMod
						module := mod.NewZyroMod(
							make(map[string]*wrapper.ZyroObject),
							publicFunctions,
							privateFunctions,
							currentFunctionName,
							unit.GetFile(),
							currentModVars,
							currentFunctionPublic,
							currentFunctionTrace,
							genericParams,
						)

						result.AddModule(unit.GetFile(), currentFunctionName, &module, unit)
						inMod = false
						expectingFun = true

						resetFlags(
							&currentFunctionName,
							&currentFunctionReturnType,
							&currentFunctionParameters,
							&currentFunctionPublic,
							&currentParameter,
							&currentFunctionBody,
						)

						genericTokens = nil
						genericTokens = make([]token.Token, 0)
						inGeneric = false
						genericDepth = 0

						currentModVars = nil
						currentModVars = make([][]token.Token, 0)
						continue
					}

					inMod = false
					expectingFun = true

					// Create the function
					function := code.NewFunction(
						currentFunctionReturnType,
						currentFunctionParameters,
						currentFunctionBody,
						currentFunctionPublic,
						strings.HasPrefix(unit.GetFile(), stdPath),
						currentFunctionTrace,
					)

					result.AddFunction(
						unit,
						unit.GetFile(),
						currentFunctionName,
						function,
					)

					resetFlags(
						&currentFunctionName,
						&currentFunctionReturnType,
						&currentFunctionParameters,
						&currentFunctionPublic,
						&currentParameter,
						&currentFunctionBody,
					)
					continue
				}
			}

			currentFunctionBody = append(currentFunctionBody, unit)
		}
	}

	return &result
}

// resetFlags resets the flags used to keep track of the parser's state
func resetFlags(
	currentFunctionName *string,
	currentFunctionReturnType *[]token.Token,
	currentFunctionParameters *map[string][]token.Token,
	currentFunctionPublic *bool,
	currentParameter *[]token.Token,
	currentFunctionBody *[]token.Token,
) {
	// Reset the current function's metadata
	*currentFunctionName = ""
	*currentFunctionReturnType = nil
	*currentFunctionReturnType = []token.Token{}

	*currentFunctionParameters = make(map[string][]token.Token)
	*currentFunctionPublic = false
	*currentParameter = nil
	*currentParameter = make([]token.Token, 0)

	*currentFunctionBody = nil
	*currentFunctionBody = make([]token.Token, 0)
}
