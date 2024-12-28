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
// and returns the inferred object
// and a boolean indicating if it is a constant value
func ToObj(
	unit token.Token,
	variables *stack.Stack,
) (wrapper.FluentObject, bool) {
	tokenType := unit.GetType()

	switch tokenType {
	case token.BoolLiteral:
		return wrapper.NewFluentObject(dummyBoolType, ""), true
	case token.StringLiteral:
		return wrapper.NewFluentObject(dummyStringType, ""), true
	case token.NumLiteral:
		return wrapper.NewFluentObject(dummyIntType, ""), true
	case token.DecimalLiteral:
		return wrapper.NewFluentObject(dummyDecimalType, ""), true
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

		return variable.GetValue(), variable.IsConstant()
	default:
		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
			"Got "+unit.GetValue(),
		)

		return wrapper.NewFluentObject(dummyNothingType, ""), true
	}
}
