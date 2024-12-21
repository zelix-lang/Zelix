package tokenUtil

import (
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
)

// FromRawType converts a raw token to a SurfObject
func FromRawType(
	token code.Token,
	variables *stack.Stack,
) object.SurfObject {
	tokenType := token.GetType()

	switch tokenType {
	case code.Bool:
		return object.NewSurfObject(object.BooleanType, "")
	case code.String:
		return object.NewSurfObject(object.StringType, "")
	case code.Num:
		return object.NewSurfObject(object.IntType, "")
	case code.Dec:
		return object.NewSurfObject(object.DecimalType, "")
		/*case code.Identifier:
		  todo!
		*/
	default:
		logger.TokenError(
			token,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return object.NewSurfObject(object.NothingType, "")
	}
}

// ToObj converts a token to a SurfObject
func ToObj(
	token code.Token,
	variables *stack.Stack,
) object.SurfObject {
	tokenType := token.GetType()

	switch tokenType {
	case code.BoolLiteral:
		return object.NewSurfObject(object.BooleanType, "")
	case code.StringLiteral:
		return object.NewSurfObject(object.StringType, "")
	case code.NumLiteral:
		return object.NewSurfObject(object.IntType, "")
	case code.DecimalLiteral:
		return object.NewSurfObject(object.DecimalType, "")
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

		return object.NewSurfObject(object.NothingType, "")
	}
}
