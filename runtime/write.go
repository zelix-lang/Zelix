package runtime

import "surf/object"

// Write writes the given string to the standard output
// without a newline character
func Write(objects ...object.SurfObject) {
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
		default:
			panic("unhandled default case")
		}
	}
}

// Writeln writes the given string to the standard output
// with a newline character
func Writeln(objects ...object.SurfObject) {
	Write(objects...)
	println()
}
