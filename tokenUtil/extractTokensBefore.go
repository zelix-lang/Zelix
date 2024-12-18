package tokenUtil

import (
	"surf/code"
	"surf/logger"
)

// ExtractTokensBefore returns the tokens before the first occurrence of the delimiter
func ExtractTokensBefore(
	tokens []code.Token,
	delimiter code.TokenType,
	handleNested bool,
	nestedStartDelimiter code.TokenType,
	nestedEndDelimiter code.TokenType,
) []code.Token {
	// Used to know if the delimiter was found
	metDelimiter := false

	blockDepth := 0
	result := make([]code.Token, 0)

	for _, token := range tokens {
		if handleNested {
			if token.GetType() == nestedStartDelimiter {
				blockDepth++
			}
			if token.GetType() == nestedEndDelimiter {
				blockDepth--

				if blockDepth < 0 {
					logger.TokenError(
						token,
						"Unmatched delimiter",
						"Match the delimiters",
					)
				}
			}
		}

		if token.GetType() == delimiter && blockDepth == 0 {
			metDelimiter = true
			break
		}

		result = append(result, token)
	}

	if !metDelimiter {
		logger.TokenError(
			tokens[len(tokens)-1],
			"Expected a delimiter",
			"Add a delimiter to the end of the list",
		)
	}

	return result
}
