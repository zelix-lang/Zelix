package structure

import (
	"fluent/code/types"
	"fluent/code/wrapper"
	"strings"
)

// MarshalType marshals the given TypeWrapper into the IR string
func MarshalType(typeWrapper wrapper.TypeWrapper, builder *strings.Builder) {
	if typeWrapper.GetType() == types.ModType {
		builder.WriteString("module")
		return
	}

	builder.WriteString(typeWrapper.Marshal())
}
