package _type

import (
	"strconv"
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
)

// TranslateType translates a token to a Surf object
func TranslateType(
	token code.Token,
	variables *stack.Stack,
) object.SurfObject {
	tokenType := token.GetType()
	value := token.GetValue()

	switch tokenType {
	case code.StringLiteral:
		return object.NewSurfObject(object.StringType, value)
	case code.BoolLiteral:
		return object.NewSurfObject(object.BooleanType, value == "true")
	case code.DecimalLiteral:
		floatValue, err := strconv.ParseFloat(value, 64)

		if err != nil {
			logger.TokenError(
				token,
				"Invalid decimal value",
				"Use a valid float value",
			)
		}

		return object.NewSurfObject(object.DecimalType, floatValue)
	case code.NumLiteral:
		intValue, err := strconv.Atoi(value)

		if err != nil {
			logger.TokenError(
				token,
				"Invalid number value",
				"Use a valid integer value",
			)
		}

		return object.NewSurfObject(object.IntType, intValue)
	case code.Identifier:
		objectValue, found := variables.Load(value)

		if !found {
			logger.TokenError(
				token,
				"Undefined reference to variable "+value,
				"The variable "+value+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		return objectValue
	default:
		logger.TokenError(
			token,
			"Invalid type",
			"Use a valid type",
		)

		return object.NewSurfObject(object.StringType, nil)
	}
}
