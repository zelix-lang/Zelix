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
	throwIfNotFound bool,
) []token.Token {
	// Used to know if the delimiter was found
	metDelimiter := false

	blockDepth := 0
	result := make([]token.Token, 0)

	for _, unit := range tokens {
		if handleNested {
			if unit.GetType() == nestedStartDelimiter {
				blockDepth++
			}
			if unit.GetType() == nestedEndDelimiter {
				blockDepth--

				if blockDepth < 0 {
					logger.TokenError(
						unit,
						"Unmatched delimiter",
						"Match the delimiters",
					)
				}
			}
		}

		if unit.GetType() == delimiter && blockDepth == 0 {
			metDelimiter = true
			break
		}

		result = append(result, unit)
	}

	if throwIfNotFound && !metDelimiter {
		logger.TokenError(
			tokens[len(tokens)-1],
			"Expected a "+TokenTypeToString(delimiter),
			"Add a delimiter to the end of the list",
		)
	}

	return result

}
