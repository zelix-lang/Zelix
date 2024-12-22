package tokenUtil

import (
	"surf/code"
	"surf/core/stack"
	"surf/logger"
	"surf/object"
	"surf/token"
)

// FromRawType converts a raw token to a SurfObject
func FromRawType(
	unit token.Token,
	mods *map[string]map[string]*code.SurfMod,
) object.SurfObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.Bool:
		return object.NewSurfObject(object.BooleanType, "")
	case token.String:
		return object.NewSurfObject(object.StringType, "")
	case token.Num:
		return object.NewSurfObject(object.IntType, "")
	case token.Dec:
		return object.NewSurfObject(object.DecimalType, "")
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

		return object.NewSurfObject(object.ModType, mod)
	default:
		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return object.NewSurfObject(object.NothingType, "")
	}
}

// ToObj converts a token to a SurfObject
func ToObj(
	unit token.Token,
	variables *stack.Stack,
) object.SurfObject {
	tokenType := unit.GetType()

	switch tokenType {
	case token.BoolLiteral:
		return object.NewSurfObject(object.BooleanType, "")
	case token.StringLiteral:
		return object.NewSurfObject(object.StringType, "")
	case token.NumLiteral:
		return object.NewSurfObject(object.IntType, "")
	case token.DecimalLiteral:
		return object.NewSurfObject(object.DecimalType, "")
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

		return object.NewSurfObject(object.NothingType, "")
	}
}
