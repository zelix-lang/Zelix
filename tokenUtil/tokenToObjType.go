package tokenUtil

import (
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
)

// FromRawType converts a raw type to a SurfObjectType
func FromRawType(
	token code.Token,
	variables *stack.StaticStack,
) object.SurfObjectType {
	tokenType := token.GetType()

	switch tokenType {
	case code.Bool:
		return object.BooleanType
	case code.String:
		return object.StringType
	case code.Num:
		return object.IntType
	case code.Dec:
		return object.DecimalType
	case code.Nothing:
		return object.NothingType
	default:
		logger.TokenError(
			token,
			"Unexpected token",
			"Expected a raw type",
		)

		return object.NothingType
	}
}

// ToObjType converts a token type to a SurfObjectType
func ToObjType(
	token code.Token,
	variables *stack.StaticStack,
) object.SurfObjectType {
	tokenType := token.GetType()

	switch tokenType {
	case code.BoolLiteral:
		return object.BooleanType
	case code.StringLiteral:
		return object.StringType
	case code.NumLiteral:
		return object.IntType
	case code.DecimalLiteral:
		return object.DecimalType
	case code.Identifier:
		variable, found := variables.Load(token.GetValue())

		if !found {
			logger.TokenError(
				token,
				"Undefined reference to variable "+token.GetValue(),
				"The variable "+token.GetValue()+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		return variable
	default:
		logger.TokenError(
			token,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return object.NothingType
	}
}
