package runtime

import (
	"zyro/code/mod"
	"zyro/code/wrapper"
)

// Write writes the given string to the standard output
// without a newline character
func Write(objects ...wrapper.ZyroObject) {
	for _, obj := range objects {
		typeWrapper := obj.GetType()
		_type := typeWrapper.GetType()

		switch _type {
		case wrapper.StringType:
			print(obj.GetValue().(string))
		case wrapper.DecimalType:
			print(obj.GetValue().(float64))
		case wrapper.IntType:
			print(obj.GetValue().(int))
		case wrapper.BooleanType:
			print(obj.GetValue().(bool))
		case wrapper.NothingType:
			print("@Zyro<Nothing>")
		default:
			mod := obj.GetValue().(mod.ZyroMod)
			print(mod.GetName())
		}
	}
}

// Writeln writes the given string to the standard output
// with a newline character
func Writeln(objects ...wrapper.ZyroObject) {
	Write(objects...)
	println()
}
