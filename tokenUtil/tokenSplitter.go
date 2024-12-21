package tokenUtil

import (
	"surf/logger"
	"surf/token"
)

// SplitTokens splits the given tokens
// using the given delimiter
func SplitTokens(
	tokens []token.Token,
	delimiter token.Type,
	nestedStartDelimiter token.Type,
	nestedEndDelimiter token.Type,
) [][]token.Token {
	result := make([][]token.Token, 0)
	current := make([]token.Token, 0)

	blockDepth := 0

	for _, unit := range tokens {
		tokenType := unit.GetType()

		if tokenType == nestedStartDelimiter {
			blockDepth++
		} else if tokenType == nestedEndDelimiter {
			blockDepth--

			if blockDepth < 0 {
				logger.TokenError(
					unit,
					"Unmatched delimiter",
					"Match the delimiters",
				)
			}
		}

		if tokenType == delimiter && blockDepth == 0 {
			if len(current) == 0 {
				logger.TokenError(
					unit,
					"Empty statement",
					"Add a statement to the list",
				)
			}

			result = append(result, current)
			current = make([]token.Token, 0)
			continue
		}

		current = append(current, unit)
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}
