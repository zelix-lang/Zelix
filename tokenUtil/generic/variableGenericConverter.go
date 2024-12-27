package generic

import (
	"fluent/code/wrapper"
	"fluent/token"
)

// ConvertVariableGenerics removes generic types from variable
// declarations
func ConvertVariableGenerics(
	declaration []token.Token,
	types map[string]wrapper.TypeWrapper,
) []token.Token {
	newBody := make([]token.Token, 0)
	inDataType := false
	inValue := false
	dataType := make([]token.Token, 0)
	value := make([]token.Token, 0)

	for _, d := range declaration {
		dType := d.GetType()
		if inDataType {
			if dType == token.Assign {
				inDataType = false
				inValue = true
			}

			dataType = append(dataType, d)
			continue
		}

		if inValue {
			value = append(value, d)
			continue
		}

		inDataType = dType == token.Colon
		newBody = append(newBody, d)
	}

	// Exclude the last token (the assign)
	usableDataType := dataType[:len(dataType)-1]
	typeWrapper := wrapper.NewTypeWrapper(usableDataType, usableDataType[0])
	newType := ConvertGeneric(typeWrapper, types)

	newBody = append(newBody, newType.MarshalTokens(usableDataType[0])...)
	// Append the assign
	newBody = append(newBody, dataType[len(dataType)-1])

	// Append the value
	newBody = append(newBody, value...)

	return newBody
}
