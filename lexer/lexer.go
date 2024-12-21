package lexer

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"surf/logger"
	"surf/token"
	"surf/util"
)

// Punctuation characters are meant to be separated tokens
// i.e.: "my_var.create" => ["my_var", ".", "create"]
// therefore, create a map with all the punctuation characters
// not a slice because that will give O(n) time complexity
// compared to O(1) time complexity with a map
var punctuation = map[rune]struct{}{
	';': {}, ',': {}, '(': {}, ')': {}, '{': {}, '}': {}, ':': {}, '+': {}, '-': {},
	'>': {}, '<': {}, '%': {}, '*': {}, '.': {}, '=': {}, '!': {}, '/': {}, '&': {},
	'|': {}, '^': {}, '[': {}, ']': {},
}

// Regex to match identifiers
var identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// Regex to match numbers
var numberRegex = regexp.MustCompile(`^[0-9]+$`)

// The standard library's path
var stdPath = os.Getenv("SURF_STANDARD_PATH")

func extractImportPath(
	result []token.Token,
	i int,
) (string, bool) {
	importPathRaw := result[i].GetValue()
	isStd := false

	if strings.HasPrefix(importPathRaw, "@std") {
		// Replace the import path with the standard library path
		// and add the .surf extension
		importPathRaw = strings.Replace(importPathRaw, "@std", stdPath, 1) + ".surf"
		isStd = true
	}

	// Get the directory of the file that the import came from
	importDir := util.DirName(result[i-1].GetFile())

	// Process the import path
	importPath := importPathRaw

	if !isStd {
		importPath = filepath.Join(importDir, importPathRaw)
	}

	return importPath, isStd
}

// validateImport validates import statements
func validateImport(
	tokens []token.Token,
	startAt int,
) {
	// Get the tokens to be validated
	importPathToken := tokens[startAt]
	semicolonToken := tokens[startAt+1]

	if importPathToken.GetType() != token.StringLiteral {
		logger.TokenError(
			importPathToken,
			"Expected a string literal as an import path",
			"Imports must be string literals",
		)
	}

	if semicolonToken.GetType() != token.Semicolon {
		logger.TokenError(
			semicolonToken,
			"Expected a semicolon after the import path",
			"Imports must end with a semicolon",
		)
	}

	// Validate the import path
	importPath, _ := extractImportPath(tokens, startAt)

	// Check if the path exists
	if !util.FileExists(importPath) {
		logger.TokenError(
			importPathToken,
			"Import path does not exist",
			"Make sure the file exists",
			"Imported path: "+importPath,
		)
	}

	// Check if the path is a directory
	if util.IsDir(importPath) {
		logger.TokenError(
			importPathToken,
			"Import path is a directory",
			"Make sure the path is a file",
		)
	}
}

// pushToken pushes the current token to the tokens array
// and clears the current token
// in O(1) time
func pushToken(
	currentToken *strings.Builder,
	tokens *[]token.Token,
	line int,
	column int,
	file string,
	input string,
	currentIndex int,
	decimalLiteral bool,
) {
	// Ignore empty tokens
	if (strings.Trim(currentToken.String(), " ")) == "" {
		return
	}

	// See if the current token is a known token
	knownToken, ok := GetKnownToken(currentToken.String())
	tokenType := token.Unknown

	if ok {
		tokenType = knownToken
	} else if decimalLiteral {
		tokenType = token.DecimalLiteral
	} else if numberRegex.MatchString(currentToken.String()) {
		tokenType = token.NumLiteral
	} else {
		// Check if the token is an identifier
		if identifierRegex.MatchString(currentToken.String()) {
			tokenType = token.Identifier
		}
	}

	*tokens = append(
		*tokens,
		*token.NewToken(
			tokenType,
			currentToken.String(),
			file,
			line,
			column,
			input,
			currentIndex,
			*currentToken,
		),
	)

	// Clear the current token
	currentToken.Reset()
}

// containsImports
// Checks if the tokens array contains any imports
// in O(n) time
func containsImports(unit token.Token) bool {
	return unit.GetType() == token.Import
}

// Lex converts a string of source dode into
// an array of tokens, processing imports
// in O(n) time
func Lex(input string, file string) []token.Token {
	// Lex the input
	result := lexSingleFile(input, file)

	// Save the seen imports in a slice
	// to detect circular imports
	var seenImports []string

	// Process imports
	for util.AnyMatch(result, containsImports) {
		for i, unit := range result {
			if unit.GetType() != token.Import {
				continue
			}

			// Validate the import
			validateImport(result, i+1)

			// Get the import path
			importPath, isStd := extractImportPath(result, i+1)

			// Remove the import statement
			result = append(result[:i], result[i+3:]...)

			// Check if we already have seen this import
			if util.AnyMatch(seenImports, func(s string) bool {
				return s == importPath
			}) {
				// Standard imports can't cause circular imports, skip
				if isStd {
					continue
				}

				// Build the import chain to print
				var importChain []string

				// Add the message to the chain so it also gets printed
				importChain = append(importChain, "File "+unit.GetFile()+" depends on its own: ")

				for spaces, _import := range seenImports {
					importChain = append(
						importChain,
						strings.Repeat(" ", spaces)+_import,
					)
				}

				logger.TokenError(
					unit,
					"Circular import detected",
					importChain...,
				)
			}

			// Read the file contents
			importContents, err := os.ReadFile(importPath)

			if err != nil {
				logger.TokenError(
					result[i],
					"Could not read the imported file",
					"Make sure the file exists",
					"File imported: "+importPath,
				)
			}

			// Lex the import contents
			importTokens := lexSingleFile(string(importContents), importPath)

			// Insert the import tokens to the end of the result
			result = append(result, importTokens...)

			if !isStd {
				seenImports = append(seenImports, importPath)
			}
		}
	}

	return result
}

// lexSingleFile converts a string of source dode into
// an array of tokens without processing imports
// in O(n) time
func lexSingleFile(input string, file string) []token.Token {
	// Used to keep track of the line and column
	line := 1
	column := 0

	// Used to keep track of the state of the lexer
	inString := false
	inStringEscape := false
	inComment := false
	inBlockComment := false

	// Used to know if the current token is a decimal literal
	decimalLiteral := false

	var result []token.Token

	// Use a builder for efficiency
	var currentToken strings.Builder

	// The input length
	inputLength := len(input)

	// Used to skip indexes
	countToIndex := 0

	// Iterate over each character
	for i, char := range input {
		if i < countToIndex {
			continue
		}

		column++

		if char == '\n' {
			line++
			column = 0

			inComment = false

			continue
		}

		if !inString && !inComment && !inBlockComment && i+1 < inputLength {
			if char == '/' && input[i+1] == '/' {
				// Push any remaining token
				pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)

				decimalLiteral = false
				inComment = true
				continue
			}

			if char == '/' && input[i+1] == '*' {
				// Push any remaining token
				pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)

				decimalLiteral = false
				inBlockComment = true
				continue
			}
		}

		// Handle exiting block comments
		if char == '/' && i-1 >= 0 && input[i-1] == '*' {
			inBlockComment = false
			continue
		}

		if inComment || inBlockComment {
			continue
		}

		if char == '"' {
			if inString {
				result = append(
					result,
					*token.NewToken(
						token.StringLiteral,
						currentToken.String(),
						file,
						line,
						column,
						input,
						i,
						currentToken,
					),
				)

				inString = false

				// Clear the current token
				currentToken.Reset()
				continue
			}

			// Push any remaining token
			pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)

			decimalLiteral = false
			inString = true
			continue
		}

		if inStringEscape {
			// No need to create a strings.Builder here
			// performance tradeoff is not worth it for a single character
			// Construct the escape sequence
			combined := "\\" + string(char)
			escaped, err := strconv.Unquote(`"` + combined + `"`)

			if err != nil {
				logger.Error("Invalid escape sequence: " + combined)

				trace, indicator := util.BuildTrace(currentToken, i, input)
				logger.Log("Full context:", trace, indicator)
				logger.Help("Use a valid escape sequence")

				os.Exit(1)
			}

			currentToken.WriteString(escaped)
			inStringEscape = false
			continue
		}

		if inString {
			if char == '\\' {
				inStringEscape = true
				// Wait for the next character to dynamically
				// construct the escape sequence
				continue
			}

			currentToken.WriteRune(char)
			continue
		}

		if char == ' ' {
			pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)
			decimalLiteral = false
			continue
		}

		// Handle arrows (->)
		if char == '>' && i-1 >= 0 && input[i-1] == '-' {
			// Remove the last token
			result = result[:len(result)-1]

			arrow := "->"

			// Create a new strings.Builder because pushToken takes a pointer
			var arrowBuilder strings.Builder

			arrowBuilder.WriteString(arrow)
			pushToken(&arrowBuilder, &result, line, column, file, input, i, decimalLiteral)
			decimalLiteral = false
			continue
		}

		// Handle increments and decrements
		if (char == '+' || char == '-') && (i+1 < inputLength && rune(input[i-1]) == char) {
			// Remove the last token
			result = result[:len(result)-1]

			incDec := string(char) + string(char)
			// Create a new strings.Builder because pushToken takes a pointer
			var incDecBuilder strings.Builder

			incDecBuilder.WriteString(incDec)
			pushToken(&incDecBuilder, &result, line, column, file, input, i, decimalLiteral)
			decimalLiteral = false

			continue
		}

		if char == '.' {
			// Handle number literals
			if currentToken.Len() > 0 && numberRegex.MatchString(currentToken.String()) {
				decimalLiteral = true
				currentToken.WriteRune(char)
				continue
			}
		}

		// Check for punctuation characters
		if _, exists := punctuation[char]; exists {
			// Push any remaining token
			pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)
			decimalLiteral = false

			// Add the punctuation character as a token
			// currentToken is already cleared, just add the punctuation character
			currentToken.WriteRune(char)
			pushToken(&currentToken, &result, line, column, file, input, i, decimalLiteral)

			continue
		}

		currentToken.WriteRune(char)
	}

	// Push the last token
	pushToken(&currentToken, &result, line, column, file, input, inputLength, decimalLiteral)

	return result
}
