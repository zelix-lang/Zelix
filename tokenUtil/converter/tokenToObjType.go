package converter

import (
	"zyro/code/types"
	"zyro/code/wrapper"
	"zyro/logger"
	"zyro/stack"
	"zyro/token"
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

// ToObj converts a token to a ZyroObject
func ToObj(
	unit token.Token,
	variables *stack.Stack,
) wrapper.ZyroObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.BoolLiteral:
		return wrapper.NewZyroObject(dummyBoolType, "")
	case token.StringLiteral:
		return wrapper.NewZyroObject(dummyStringType, "")
	case token.NumLiteral:
		return wrapper.NewZyroObject(dummyIntType, "")
	case token.DecimalLiteral:
		return wrapper.NewZyroObject(dummyDecimalType, "")
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

		return wrapper.NewZyroObject(dummyNothingType, "")
	}
}
