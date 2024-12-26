package converter

import (
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/stack"
	"fluent/token"
)

var dummyBoolType = wrapper.ForceNewTypeWrapper(
	"bool",
	[]wrapper.TypeWrapper{},
	types.BooleanType,
)

var dummyNothingType = wrapper.ForceNewTypeWrapper(
	"nothing",
	[]wrapper.TypeWrapper{},
	types.NothingType,
)

var dummyStringType = wrapper.ForceNewTypeWrapper(
	"str",
	[]wrapper.TypeWrapper{},
	types.StringType,
)

var dummyIntType = wrapper.ForceNewTypeWrapper(
	"num",
	[]wrapper.TypeWrapper{},
	types.IntType,
)

var dummyDecimalType = wrapper.ForceNewTypeWrapper(
	"dec",
	[]wrapper.TypeWrapper{},
	types.DecimalType,
)

// ToObj converts a token to a FluentObject
func ToObj(
	unit token.Token,
	variables *stack.Stack,
) wrapper.FluentObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.BoolLiteral:
		return wrapper.NewFluentObject(dummyBoolType, "")
	case token.StringLiteral:
		return wrapper.NewFluentObject(dummyStringType, "")
	case token.NumLiteral:
		return wrapper.NewFluentObject(dummyIntType, "")
	case token.DecimalLiteral:
		return wrapper.NewFluentObject(dummyDecimalType, "")
	case token.Identifier:
		variable, found := variables.Load(unit.GetValue())

		if !found {
			logger.TokenError(
				unit,
				"Undefined reference to variable "+unit.GetValue(),
				"The variable "+unit.GetValue()+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		return variable.GetValue()
	default:
		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return wrapper.NewFluentObject(dummyNothingType, "")
	}
}
