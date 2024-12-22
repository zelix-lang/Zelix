package tokenUtil

import (
	"zyro/code"
	"zyro/core/stack"
	"zyro/logger"
	"zyro/object"
	"zyro/token"
)

// FromRawType converts a raw token to a ZyroObject
func FromRawType(
	unit token.Token,
	mods *map[string]map[string]*code.ZyroMod,
) object.ZyroObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.Bool:
		return object.NewZyroObject(object.BooleanType, "")
	case token.String:
		return object.NewZyroObject(object.StringType, "")
	case token.Num:
		return object.NewZyroObject(object.IntType, "")
	case token.Dec:
		return object.NewZyroObject(object.DecimalType, "")
	case token.Identifier:
		mod, found, _ := code.FindMod(mods, unit.GetValue(), unit.GetFile())

		if !found {
			logger.TokenError(
				unit,
				"Undefined reference to module "+unit.GetValue(),
				"The module "+unit.GetValue()+" was not found in the current scope",
				"Add the variable to the current scope",
			)
		}

		return object.NewZyroObject(object.ModType, mod)
	default:
		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return object.NewZyroObject(object.NothingType, "")
	}
}

// ToObj converts a token to a ZyroObject
func ToObj(
	unit token.Token,
	variables *stack.Stack,
) object.ZyroObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.BoolLiteral:
		return object.NewZyroObject(object.BooleanType, "")
	case token.StringLiteral:
		return object.NewZyroObject(object.StringType, "")
	case token.NumLiteral:
		return object.NewZyroObject(object.IntType, "")
	case token.DecimalLiteral:
		return object.NewZyroObject(object.DecimalType, "")
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

		return object.NewZyroObject(object.NothingType, "")
	}
}
