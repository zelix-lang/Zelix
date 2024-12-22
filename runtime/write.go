package runtime

import (
	"zyro/code/mod"
	"zyro/code/types"
	"zyro/code/wrapper"
)

// Write writes the given string to the standard output
// without a newline character
func Write(objects ...wrapper.ZyroObject) {
	for _, obj := range objects {
		typeWrapper := obj.GetType()
		_type := typeWrapper.GetType()

		switch _type {
		case types.StringType:
			print(obj.GetValue().(string))
		case types.DecimalType:
			print(obj.GetValue().(float64))
		case types.IntType:
			print(obj.GetValue().(int))
		case types.BooleanType:
			print(obj.GetValue().(bool))
		case types.NothingType:
			print("@Zyro<Nothing>")
		default:
			module := obj.GetValue().(mod.ZyroMod)
			print(module.GetName())
		}
	}
}

// Writeln writes the given string to the standard output
// with a newline character
func Writeln(objects ...wrapper.ZyroObject) {
	Write(objects...)
	println()
}
