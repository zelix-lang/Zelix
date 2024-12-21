package _type

import (
	"strconv"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
)

// TranslateType translates a token to a Surf object
func TranslateType(
	unit token.Token,
	variables *stack.Stack,
) object.SurfObject {
	tokenType := unit.GetType()
	value := unit.GetValue()

	switch tokenType {
	case token.StringLiteral:
		return object.NewSurfObject(object.StringType, value)
	case token.BoolLiteral:
		return object.NewSurfObject(object.BooleanType, value == "true")
	case token.DecimalLiteral:
		floatValue, err := strconv.ParseFloat(value, 64)

		if err != nil {
			logger.TokenError(
				unit,
				"Invalid decimal value",
				"Use a valid float value",
			)
		}

		return object.NewSurfObject(object.DecimalType, floatValue)
	case token.NumLiteral:
		intValue, err := strconv.Atoi(value)

		if err != nil {
			logger.TokenError(
				unit,
				"Invalid number value",
				"Use a valid integer value",
			)
		}

		return object.NewSurfObject(object.IntType, intValue)
	case token.Identifier:
		objectValue, found := variables.Load(value)

		if !found {
			logger.TokenError(
				unit,
				"Undefined reference to variable "+value,
				"The variable "+value+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		return objectValue.GetValue()
	default:
		logger.TokenError(
			unit,
			"Invalid type",
			"Use a valid type",
		)

		return object.NewSurfObject(object.StringType, nil)
	}
}
