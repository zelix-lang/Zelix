package wrapper

import (
	"surf/logger"
	"surf/token"
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
