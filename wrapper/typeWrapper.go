package wrapper

import (
	"zyro/logger"
	"zyro/token"
	"zyro/tokenUtil"
)

// TypeWrapper is a wrapper for a data type
// based on the given tokens
// For example: Result<str, bool>
type TypeWrapper struct {
	baseType   string        // Result
	parameters []TypeWrapper // [str, bool]
}

func NewTypeWrapper(tokens []token.Token, trace token.Token) TypeWrapper {
	// Parse the tokens
	if len(tokens) == 0 {
		logger.TokenError(
			trace,
			"Invalid data type",
			"This data type is empty",
		)
	}

	baseType := tokens[0]
	parameters := make([]TypeWrapper, 0)

	// Check if the data type has parameters
	if len(tokens) > 1 {
		if tokens[1].GetType() != token.LessThan {
			logger.TokenError(
				tokens[1],
				"Invalid data type",
				"Expected '<' after the base type",
			)
		}

		if tokens[len(tokens)-1].GetType() != token.GreaterThan {
			logger.TokenError(
				tokens[len(tokens)-1],
				"Invalid data type",
				"Expected '>' at the end of the data type",
			)
		}
	}

	// Split by commas what's inside the '<' and '>'
	// For example: Result<str, bool>
	paramsTokens := tokenUtil.SplitTokens(
		tokens[2:len(tokens)-1],
		token.Comma,
		token.LessThan,
		token.GreaterThan,
	)
	for _, paramTokens := range paramsTokens {
		parameters = append(parameters, NewTypeWrapper(paramTokens, trace))
	}

	return TypeWrapper{
		baseType:   baseType.GetValue(),
		parameters: parameters,
	}
}

func (tw *TypeWrapper) GetBaseType() string {
	return tw.baseType
}

func (tw *TypeWrapper) GetParameters() []TypeWrapper {
	return tw.parameters
}
