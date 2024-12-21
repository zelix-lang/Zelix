package tokenUtil

import (
	"surf/logger"
	"surf/token"
)

// ExtractTokensBefore returns the tokens before the first occurrence of the delimiter
func ExtractTokensBefore(
	tokens []token.Token,
	delimiter token.Type,
	handleNested bool,
	nestedStartDelimiter token.Type,
	nestedEndDelimiter token.Type,
) []token.Token {
	// Used to know if the delimiter was found
	metDelimiter := false

	blockDepth := 0
	result := make([]token.Token, 0)

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
			"Expected a "+TokenTypeToString(delimiter),
			"Add a delimiter to the end of the list",
		)
	}

	return result

}
