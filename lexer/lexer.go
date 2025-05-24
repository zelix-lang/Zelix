/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
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
	"strconv"
	"strings"
	"unicode"
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

// pushToken pushes the current token to the tokens array
// and clears the current token
// in O(1) time
func pushToken(
	currentToken *strings.Builder,
	tokens *[]token.Token,
	line int,
	column int,
	file string,
	isIdentifier bool,
	isNumber bool,
	decimalLiteral bool,
) {
	value := currentToken.String()
	// Ignore empty tokens
	if len(value) == 0 {
		return
	}

	// See if the current token is a known token
	knownToken, ok := getKnownToken(value)
	newToken := token.NewToken(
		token.Unknown,
		nil,
		file,
		line,
		column,
	)

	if ok {
		newToken.TokenType = knownToken
	} else if decimalLiteral {
		newToken.TokenType = token.DecimalLiteral
		newToken.Value = &value
	} else if isNumber {
		newToken.TokenType = token.NumLiteral
		newToken.Value = &value
	} else if isIdentifier {
		newToken.TokenType = token.Identifier
		newToken.Value = &value
	}

	*tokens = append(
		*tokens,
		newToken,
	)

	// Clear the current token
	currentToken.Reset()
}

// Lex converts a string of source dode into
// an array of tokens without processing imports
// in O(n) time
func Lex(input string, file string) ([]token.Token, *Error) {
	// Used to keep track of the line and column
	line := 1
	column := 0

	// Used to keep track of the state of the lexer
	inString := false
	inStringEscape := false
	inComment := false
	inBlockComment := false

	// Used to know if the current token is a decimal literal
	isIdentifier := true
	tokenIdx := 0
	isNumber := false
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

		// Ignore carriage returns on Windows
		if char == '\r' {
			continue
		}

		if char == '\n' {
			if inString {
				return make([]token.Token, 0), &Error{
					Line:    line,
					Column:  column,
					File:    file,
					Message: "String not closed",
				}
			}

			line++
			column = 0
			inComment = false

			continue
		}

		if !inString && !inComment && !inBlockComment && i+1 < inputLength {
			if char == '/' && input[i+1] == '/' {
				// Push any remaining token
				pushToken(
					&currentToken,
					&result,
					line,
					column,
					file,
					isIdentifier,
					isNumber,
					decimalLiteral,
				)

				decimalLiteral = false
				inComment = true
				isIdentifier = false
				tokenIdx = 0
				continue
			}

			if char == '/' && input[i+1] == '*' {
				// Push any remaining token
				pushToken(
					&currentToken,
					&result,
					line,
					column,
					file,
					isIdentifier,
					isNumber,
					decimalLiteral,
				)

				decimalLiteral = false
				inBlockComment = true
				isIdentifier = false
				tokenIdx = 0
				continue
			}
		}

		// Handle exiting block comments
		if char == '/' && i-1 >= 0 && input[i-1] == '*' {
			inBlockComment = false
			isIdentifier = true
			tokenIdx = 0
			continue
		}

		if inComment || inBlockComment {
			continue
		}

		if tokenIdx == 0 && (char == '_' || unicode.IsLetter(char)) {
			isIdentifier = true
			isNumber = false
			tokenIdx++
		} else if tokenIdx == 0 && unicode.IsDigit(char) {
			isIdentifier = false
			isNumber = true
			tokenIdx++
		}

		if char == '"' {
			if inString {
				value := currentToken.String()

				result = append(
					result,
					token.NewToken(
						token.StringLiteral,
						&value,
						file,
						line,
						column,
					),
				)

				inString = false
				isIdentifier = true
				isNumber = false
				decimalLiteral = false
				tokenIdx = 0

				// Clear the current token
				currentToken.Reset()
				continue
			}

			// Push any remaining token
			pushToken(
				&currentToken,
				&result,
				line,
				column,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)

			decimalLiteral = false
			isIdentifier = false
			isNumber = false
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
				return make([]token.Token, 0), &Error{
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
			pushToken(
				&currentToken,
				&result,
				line,
				column-1,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)

			decimalLiteral = false
			isIdentifier = true
			isNumber = false
			tokenIdx = 0
			continue
		}

		// Handle arrows (->)
		if char == '-' && i+1 < inputLength && input[i+1] == '>' {
			// Skip the next token
			countToIndex = i + 2

			// Create a new strings.Builder because pushToken takes a pointer
			var arrowBuilder strings.Builder
			arrowBuilder.WriteString("->")

			pushToken(
				&arrowBuilder,
				&result,
				line,
				column,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)
			decimalLiteral = false
			isIdentifier = false
			isNumber = false
			tokenIdx = 0
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
			pushToken(
				&incDecBuilder,
				&result,
				line,
				column,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)

			decimalLiteral = false
			isIdentifier = true
			isNumber = false
			tokenIdx = 0

			continue
		}

		if char == '.' && isNumber {
			if decimalLiteral {
				// Make sure we don't have multiple decimal points
				return make([]token.Token, 0), &Error{
					Line:    line,
					Column:  column,
					File:    file,
					Message: "Invalid decimal literal",
				}
			}

			// Handle number literals
			decimalLiteral = true
			isIdentifier = false
			currentToken.WriteRune(char)
			continue
		}

		// Handle "==", "!=", ">=", "<=", "->"
		if i+1 < inputLength && chainableTokens[char] == 1 && input[i+1] == '=' {
			// Skip the tokens
			countToIndex = i + 2

			// Create a new strings.Builder because pushToken takes a pointer
			var eqBuilder strings.Builder

			eqBuilder.WriteRune(char)
			eqBuilder.WriteByte(input[i+1])
			pushToken(
				&eqBuilder,
				&result,
				line,
				column,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)

			decimalLiteral = false
			isIdentifier = true
			isNumber = false
			tokenIdx = 0
			continue
		}

		// Check for punctuation characters
		if _, exists := punctuation[char]; exists {
			// Push any remaining token
			pushToken(
				&currentToken,
				&result,
				line,
				column-1,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)

			decimalLiteral = false
			isIdentifier = false
			isNumber = false

			// Add the punctuation character as a token
			// currentToken is already cleared, just add the punctuation character
			currentToken.WriteRune(char)
			pushToken(
				&currentToken,
				&result,
				line,
				column,
				file,
				isIdentifier,
				isNumber,
				decimalLiteral,
			)
			isIdentifier = true
			tokenIdx = 0

			continue
		}

		currentToken.WriteRune(char)
	}

	if inBlockComment {
		return make([]token.Token, 0), &Error{
			Line:    line,
			Column:  column,
			File:    file,
			Message: "Block comment not closed",
		}
	}

	// Push the last token
	pushToken(
		&currentToken,
		&result,
		line,
		column,
		file,
		isIdentifier,
		isNumber,
		decimalLiteral,
	)

	return result, nil
}
