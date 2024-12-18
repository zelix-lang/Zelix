package tokenUtil

import (
	"surf/code"
	"surf/logger"
)

// SplitTokens splits the given tokens
// using the given delimiter
func SplitTokens(
	tokens []code.Token,
	delimiter code.TokenType,
	nestedStartDelimiter code.TokenType,
	nestedEndDelimiter code.TokenType,
) [][]code.Token {
	result := make([][]code.Token, 0)
	current := make([]code.Token, 0)

	blockDepth := 0

	for _, token := range tokens {
		tokenType := token.GetType()

		if tokenType == nestedStartDelimiter {
			blockDepth++
		} else if tokenType == nestedEndDelimiter {
			blockDepth--

			if blockDepth < 0 {
				logger.TokenError(
					token,
					"Unmatched delimiter",
					"Match the delimiters",
				)
			}
		}

		if tokenType == delimiter && blockDepth == 0 {
			if len(current) == 0 {
				logger.TokenError(
					token,
					"Empty statement",
					"Add a statement to the list",
				)
			}

			result = append(result, current)
			current = make([]code.Token, 0)
			continue
		}

		current = append(current, token)
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}
