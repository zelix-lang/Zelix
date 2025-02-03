/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package lexer

import (
	"fluent/token"
	"regexp"
	"strconv"
	"strings"
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

// Characters that make a combined assignment (+=, -=, *=, /=)
// Use a map for O(1) lookup
var chainableTokens = map[rune]int{
	'!': 1, '>': 1, '<': 1, '=': 1,
}

// Regex to match identifiers
var identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// Regex to match numbers
var numberRegex = regexp.MustCompile(`^[0-9]+$`)

// pushToken pushes the current token to the tokens array
// and clears the current token
// in O(1) time
func pushToken(
	currentToken *strings.Builder,
	tokens *[]token.Token,
	line int,
	column int,
	file string,
	decimalLiteral bool,
) {
	value := currentToken.String()
	// Ignore empty tokens
	if len(value) == 0 {
		return
	}

	// See if the current token is a known token
	knownToken, ok := getKnownToken(value)
	tokenType := token.Unknown

	if ok {
		tokenType = knownToken
	} else if decimalLiteral {
		tokenType = token.DecimalLiteral
	} else if numberRegex.MatchString(value) {
		tokenType = token.NumLiteral
	} else {
		// Check if the token is an identifier
		if identifierRegex.MatchString(value) {
			tokenType = token.Identifier
		}
	}

	*tokens = append(
		*tokens,
		token.NewToken(
			tokenType,
			value,
			file,
			line,
			column,
		),
	)

	// Clear the current token
	currentToken.Reset()
}

// Lex converts a string of source dode into
// an array of tokens without processing imports
// in O(n) time
func Lex(input string, file string) ([]token.Token, Error) {
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
			column++
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
				pushToken(&currentToken, &result, line, column, file, decimalLiteral)

				decimalLiteral = false
				inComment = true
				continue
			}

			if char == '/' && input[i+1] == '*' {
				// Push any remaining token
				pushToken(&currentToken, &result, line, column, file, decimalLiteral)

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
					token.NewToken(
						token.StringLiteral,
						currentToken.String(),
						file,
						line,
						column,
					),
				)

				inString = false

				// Clear the current token
				currentToken.Reset()
				continue
			}

			// Push any remaining token
			pushToken(&currentToken, &result, line, column, file, decimalLiteral)

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
				return make([]token.Token, 0), Error{
					Line:    line,
					Column:  column,
					File:    file,
					Message: "Invalid escape sequence",
				}
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
			pushToken(&currentToken, &result, line, column-1, file, decimalLiteral)
			decimalLiteral = false
			continue
		}

		// Handle arrows (->)
		if char == '-' && i+1 < inputLength && input[i+1] == '>' {
			// Skip the next token
			countToIndex = i + 2

			// Create a new strings.Builder because pushToken takes a pointer
			var arrowBuilder strings.Builder
			arrowBuilder.WriteString("->")

			pushToken(&arrowBuilder, &result, line, column, file, decimalLiteral)
			decimalLiteral = false
			continue
		}

		// Handle increments, decrements, and logical operators
		if (char == '|' || char == '&') && (i+1 < inputLength && rune(input[i+1]) == char) {
			incDec := string(char) + string(char)

			// Skip the next token
			countToIndex = i + 2

			// Create a new strings.Builder because pushToken takes a pointer
			var incDecBuilder strings.Builder

			incDecBuilder.WriteString(incDec)
			pushToken(&incDecBuilder, &result, line, column, file, decimalLiteral)
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

		// Handle "==", "!=", ">=", "<=", "->"
		if i+1 < inputLength && chainableTokens[char] == 1 && input[i+1] == '=' {
			// Skip the tokens
			countToIndex = i + 2

			// Create a new strings.Builder because pushToken takes a pointer
			var eqBuilder strings.Builder

			eqBuilder.WriteRune(char)
			eqBuilder.WriteByte(input[i+1])
			pushToken(&eqBuilder, &result, line, column, file, decimalLiteral)
			decimalLiteral = false
			continue
		}

		// Check for punctuation characters
		if _, exists := punctuation[char]; exists {
			// Push any remaining token
			pushToken(&currentToken, &result, line, column-1, file, decimalLiteral)
			decimalLiteral = false

			// Add the punctuation character as a token
			// currentToken is already cleared, just add the punctuation character
			currentToken.WriteRune(char)
			pushToken(&currentToken, &result, line, column, file, decimalLiteral)

			continue
		}

		currentToken.WriteRune(char)
	}

	if inBlockComment {
		return make([]token.Token, 0), Error{
			Line:    line,
			Column:  column,
			File:    file,
			Message: "Block comment not closed",
		}
	}

	// Push the last token
	pushToken(&currentToken, &result, line, column, file, decimalLiteral)

	return result, Error{}
}
