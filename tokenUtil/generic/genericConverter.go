package generic

import (
	"fluent/code/types"
	"fluent/code/wrapper"
)

// ConvertGeneric returns the given TypeWrapper
// without any generic parameters
func ConvertGeneric(template wrapper.TypeWrapper, buildTypes map[string]wrapper.TypeWrapper) wrapper.TypeWrapper {
	if template.GetType() != types.ModType {
		return template
	}

	params := template.GetParameters()
	newParams := make([]wrapper.TypeWrapper, len(params))

	for i, param := range params {
		newParams[i] = ConvertGeneric(param, buildTypes)
	}

	// This should be valid at all times due to previous checks
	equivalent, _ := buildTypes[template.GetBaseType()]
	return equivalent
}
