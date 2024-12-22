package runtime

import (
	"zyro/code"
	"zyro/object"
)

// Write writes the given string to the standard output
// without a newline character
func Write(objects ...object.ZyroObject) {
	for _, obj := range objects {
		_type := obj.GetType()

		switch _type {
		case object.StringType:
			print(obj.GetValue().(string))
		case object.DecimalType:
			print(obj.GetValue().(float64))
		case object.IntType:
			print(obj.GetValue().(int))
		case object.BooleanType:
			print(obj.GetValue().(bool))
		case object.NothingType:
			print("@Zyro<Nothing>")
		default:
			mod := obj.GetValue().(code.ZyroMod)
			print(mod.GetName())
		}
	}
}

// Writeln writes the given string to the standard output
// with a newline character
func Writeln(objects ...object.ZyroObject) {
	Write(objects...)
	println()
}
