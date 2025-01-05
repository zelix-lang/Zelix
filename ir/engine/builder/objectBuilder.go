package builder

import (
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/token"
)

// BuildObject builds a statement into a FluentObject
func BuildObject(statement []token.Token) wrapper.FluentObject {
	// Manage single token statements
	if len(statement) == 1 {
		unit := statement[0]
		unitType := unit.GetType()

		var finalWrapper wrapper.TypeWrapper
		switch unitType {
		case token.StringLiteral:
			finalWrapper = wrapper.ForceNewTypeWrapper(
				"str",
				make([]wrapper.TypeWrapper, 0),
				types.StringType,
			)
		case token.NumLiteral:
			finalWrapper = wrapper.ForceNewTypeWrapper(
				"num",
				make([]wrapper.TypeWrapper, 0),
				types.IntType,
			)
		case token.BoolLiteral:
			finalWrapper = wrapper.ForceNewTypeWrapper(
				"bool",
				make([]wrapper.TypeWrapper, 0),
				types.BooleanType,
			)
		case token.DecimalLiteral:
			finalWrapper = wrapper.ForceNewTypeWrapper(
				"dec",
				make([]wrapper.TypeWrapper, 0),
				types.DecimalType,
			)
		default:
			break
		}

		return wrapper.NewFluentObject(
			finalWrapper,
			unit.GetValue(),
		)
	}

	return wrapper.FluentObject{}
}
