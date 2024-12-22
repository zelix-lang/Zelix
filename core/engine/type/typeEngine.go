package _type

import (
	"strconv"
	"zyro/core/stack"
	"zyro/logger"
	"zyro/object"
	"zyro/token"
)

// TranslateType translates a token to a Zyro object
func TranslateType(
	unit token.Token,
	variables *stack.Stack,
) object.ZyroObject {
	tokenType := unit.GetType()
	value := unit.GetValue()

	switch tokenType {
	case token.StringLiteral:
		return object.NewZyroObject(object.StringType, value)
	case token.BoolLiteral:
		return object.NewZyroObject(object.BooleanType, value == "true")
	case token.DecimalLiteral:
		floatValue, err := strconv.ParseFloat(value, 64)

		if err != nil {
			logger.TokenError(
				unit,
				"Invalid decimal value",
				"Use a valid float value",
			)
		}

		return object.NewZyroObject(object.DecimalType, floatValue)
	case token.NumLiteral:
		intValue, err := strconv.Atoi(value)

		if err != nil {
			logger.TokenError(
				unit,
				"Invalid number value",
				"Use a valid integer value",
			)
		}

		return object.NewZyroObject(object.IntType, intValue)
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

		return object.NewZyroObject(object.StringType, nil)
	}
}
