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

package util

import "fluent/token"

// ExtractTokensBefore extracts all tokens from the given slice that appear
// before the specified delimiter token type. It returns a slice of these tokens.
//
// Parameters:
// - tokens: A slice of token.Token representing the tokens to be processed.
// - delimiter: A token.Type representing the delimiter token type.
// - handleNested: A boolean indicating whether to handle nested delimiters.
// - nestedDelimiter: A token.Type representing the nested delimiter token type.
// - nestedEndDelimiter: A token.Type representing the nested end delimiter token type.
// - delimiterRequired: A boolean indicating whether the delimiter is required.
//
// Returns:
// - A slice of token.Token containing all tokens before the delimiter.
func ExtractTokensBefore(
	tokens []token.Token,
	delimiter []token.Type,
	handleNested bool,
	nestedDelimiter token.Type,
	nestedEndDelimiter token.Type,
	delimiterRequired bool,
) []token.Token {
	// Convert the delimiter to a map for faster lookup
	delimiterMap := make(map[token.Type]bool)

	for _, d := range delimiter {
		delimiterMap[d] = true
	}

	// Used to keep track of the current nesting level
	nestingLevel := 0

	// Used to recognize invalid results
	hasMetDelimiter := false

	result := make([]token.Token, 0)

	for _, unit := range tokens {
		// Handle nested delimiters
		if handleNested {
			if unit.TokenType == nestedDelimiter {
				nestingLevel++
			} else if unit.TokenType == nestedEndDelimiter && nestingLevel > 0 {
				nestingLevel--

				// Handle invalid nesting levels
				if nestingLevel < 0 {
					return nil
				}
			}
		}

		if nestingLevel == 0 {
			if _, ok := delimiterMap[unit.TokenType]; ok {
				hasMetDelimiter = true
				break
			}
		}

		result = append(result, unit)
	}

	// Handle invalid nesting levels and invalid results
	if nestingLevel > 0 || (!hasMetDelimiter && delimiterRequired) {
		return nil
	}

	return result
}
