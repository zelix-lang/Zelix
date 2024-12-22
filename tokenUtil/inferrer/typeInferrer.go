package inferrer

import (
	"zyro/code/types"
	"zyro/logger"
	"zyro/token"
)

// InferFromRawType infers a token's type
func InferFromRawType(unit token.Token, allowLiterals bool) types.ZyroObjectType {
	tokenType := unit.GetType()

	if allowLiterals {
		switch tokenType {
		case token.BoolLiteral:
			return types.BooleanType
		case token.StringLiteral:
			return types.StringType
		case token.NumLiteral:
			return types.IntType
		case token.DecimalLiteral:
			return types.DecimalType
		default:
			logger.TokenError(
				unit,
				"Unexpected token",
				"Expected an identifier, a literal or a variable",
			)

			return types.NothingType
		}
	}

	switch tokenType {
	case token.Nothing:
		return types.NothingType
	case token.Bool:
		return types.BooleanType
	case token.String:
		return types.StringType
	case token.Num:
		return types.IntType
	case token.Dec:
		return types.DecimalType
	case token.Identifier:
		return types.ModType
	default:
		logger.TokenError(
			unit,
			"Unexpected token",
			"Expected an identifier, a literal or a variable",
		)

		return types.NothingType
	}
}
