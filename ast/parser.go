package ast

import (
	"os"
	"strings"
	"surf/code"
	"surf/concurrent"
	"surf/logger"
	"surf/object"
	"surf/tokenUtil"
)

// The standard library's path
var stdPath = os.Getenv("SURF_STANDARD_PATH")

// Parse parses the given tokens into a FileCode
func Parse(tokens []code.Token, allowMods bool, allowInlineVars bool) *FileCode {
	result := FileCode{
		functions: make(map[string]map[string]*Function),
		modules:   make(map[string]*SurfMod),
	}

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
	var currentFunctionReturnType []code.Token
	currentFunctionParameters := make(map[string][]code.Token)

	var currentFunctionPublic bool
	currentParameter := make([]code.Token, 0)
	currentFunctionBody := make([]code.Token, 0)

	currentModVars := make([][]code.Token, 0)

	// Used to skip tokens
	skipToIndex := 0

	for i, token := range tokens {
		if i < skipToIndex {
			continue
		}

		tokenType := token.GetType()

		if tokenType == code.Let || tokenType == code.Const {
			if !allowInlineVars && !inMod {
				logger.TokenError(
					token,
					"Inline variable declarations are not allowed here",
					"Remove the inline variable declaration",
				)
			}

			// Extract the statement
			statement := tokenUtil.ExtractTokensBefore(
				tokens[i:],
				code.Semicolon,
				// Don't handle nested statements here
				false,
				code.Unknown,
				code.Unknown,
			)

			currentModVars = append(currentModVars, statement)

			skipToIndex = i + len(statement) + 1
			continue
		}

		if tokenType == code.Pub {
			if inFunction {
				logger.TokenError(
					token,
					"Cannot declare a public function inside another function",
					"Move the function outside of the current function",
				)
			}

			// Make sure there is a next token and that it's in the same file
			if i == len(tokens)-1 || tokens[i+1].GetFile() != token.GetFile() {
				logger.TokenError(
					token,
					"Expected a function declaration after the public keyword",
					"Add a function declaration after the public keyword",
				)
			}

			if currentFunctionPublic {
				logger.TokenError(
					token,
					"Cannot declare a function as public twice",
					"Remove the extra public keyword",
				)
			}

			currentFunctionPublic = true
		} else if expectingFun {
			if tokenType == code.Mod {
				if !allowMods {
					logger.TokenError(
						token,
						"Modules are not allowed here",
						"Remove the module declaration",
					)
				}

				inMod = true
				expectingFun = false
				expectingFunName = true
				continue
			}

			if tokenType != code.Function {
				logger.TokenError(
					token,
					"Expected a function declaration",
					"Add a function declaration",
				)
			}

			inFunction = true
			expectingFun = false
			expectingFunName = true
		} else if expectingFunName {
			// Only allow identifiers as function names
			if tokenType != code.Identifier {
				logger.TokenError(
					token,
					"Expected a function name",
					"Add a function name",
				)
			}

			currentFunctionName = token.GetValue()

			if inMod {
				expectingOpenCurly = true
				expectingFunName = false
				continue
			}

			expectingFunName = false
			expectingOpenParen = true
		} else if expectingOpenCurly {
			// This case only happens for modules
			if tokenType != code.OpenCurly {
				logger.TokenError(
					token,
					"Expected an opening curly brace",
					"Add an opening curly brace",
				)
			}

			expectingOpenCurly = false
			// Go directly to parse the body of the mod
		} else if expectingOpenParen {
			if tokenType != code.OpenParen {
				logger.TokenError(
					token,
					"Expected an open parenthesis",
					"Add an open parenthesis",
				)
			}

			expectingOpenParen = false
			expectingArgs = true
		} else if expectingArgs {
			if tokenType == code.CloseParen {
				// Parse the return type
				if len(currentParameter) > 0 {
					currentFunctionParameters[currentParameterName] = currentParameter
				}

				currentParameter = nil
				currentParameter = []code.Token{}
				expectingArgs = false
				expectingArrow = true
				continue
			}

			// Take the parameter's name
			// only identifiers are allowed here
			if tokenType != code.Identifier {
				logger.TokenError(
					token,
					"Expected a parameter name",
					"Add a parameter name",
				)
			}

			currentParameterName = token.GetValue()
			expectingArgs = false
			expectingArgColon = true
		} else if expectingArgColon {
			if tokenType != code.Colon {
				logger.TokenError(
					token,
					"Expected a colon after the parameter name",
					"Add a colon after the parameter name",
				)
			}

			expectingArgColon = false
			expectingArgType = true
		} else if expectingArgType {
			if tokenType == code.Comma {
				// Parse next parameter
				if len(currentParameter) > 0 {
					currentFunctionParameters[currentParameterName] = currentParameter
				}

				currentParameter = nil
				currentParameter = []code.Token{}
				expectingArgType = false
				expectingArgs = true
				continue
			}

			if tokenType == code.CloseParen {
				// Parse the return type
				if len(currentParameter) > 0 {
					currentFunctionParameters[currentParameterName] = currentParameter
				}

				currentParameter = nil
				currentParameter = []code.Token{}
				expectingArgs = false
				expectingArrow = true
				expectingArgType = false
				continue
			}

			currentParameter = append(currentParameter, token)
		} else if expectingArrow {
			if tokenType == code.OpenCurly {
				returnType := code.NewToken(
					code.Nothing,
					"nothing",
					token.GetFile(),
					token.GetLine(),
					token.GetColumn(),
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

			if tokenType != code.Arrow {
				logger.TokenError(
					token,
					"Expected an arrow after the function parameters",
					"Add an arrow after the function parameters",
				)
			}

			expectingArrow = false
			expectingReturnType = true
		} else if expectingReturnType {
			if tokenType == code.OpenCurly {
				// End of return type
				expectingReturnType = false
				continue
			}

			currentFunctionReturnType = append(currentFunctionReturnType, token)
		} else {
			if !inFunction && !inMod {
				logger.TokenError(
					token,
					"Unexpected token",
					"Remove the token",
					"You can't have a token outside of a function",
				)
			}

			if tokenType == code.OpenCurly {
				// Start of function body
				blockDepth++
			} else if tokenType == code.CloseCurly {
				// End of function body
				blockDepth--

				if blockDepth < 0 {
					logger.TokenError(
						token,
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
						privateFunctions := make(map[string]*Function)
						publicFunctions := make(map[string]*Function)

						for _, functions := range modFunctions.functions {
							for name, function := range functions {
								if function.public {
									publicFunctions[name] = function
								} else {
									privateFunctions[name] = function
								}
							}
						}

						// Wrap the functions inside a SurfMod
						mod := NewSurfMod(
							concurrent.NewTypedConcurrentMap[string, object.SurfObject](),
							publicFunctions,
							privateFunctions,
							currentFunctionName,
							token.GetFile(),
							currentModVars,
						)

						result.modules[currentFunctionName] = &mod
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

						continue
					}

					inMod = false
					expectingFun = true

					// Create the function
					function := Function{
						returnType: currentFunctionReturnType,
						parameters: currentFunctionParameters,
						body:       currentFunctionBody,
						public:     currentFunctionPublic,
						std:        strings.HasPrefix(token.GetFile(), stdPath),
						trace:      token,
					}

					result.AddFunction(
						token,
						token.GetFile(),
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

			currentFunctionBody = append(currentFunctionBody, token)
		}
	}

	return &result
}

// resetFlags resets the flags used to keep track of the parser's state
func resetFlags(
	currentFunctionName *string,
	currentFunctionReturnType *[]code.Token,
	currentFunctionParameters *map[string][]code.Token,
	currentFunctionPublic *bool,
	currentParameter *[]code.Token,
	currentFunctionBody *[]code.Token,
) {
	// Reset the current function's metadata
	*currentFunctionName = ""
	*currentFunctionReturnType = nil
	*currentFunctionReturnType = []code.Token{}

	*currentFunctionParameters = make(map[string][]code.Token)
	*currentFunctionPublic = false
	*currentParameter = nil
	*currentParameter = make([]code.Token, 0)

	*currentFunctionBody = nil
	*currentFunctionBody = make([]code.Token, 0)
}
