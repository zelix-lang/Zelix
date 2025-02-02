/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package variable

import "fluent/analyzer/object"

// Variable represents a variable in the Fluent programming language.
type Variable struct {
	// Constant indicates if the variable is a constant.
	Constant bool
	// Value holds the value of the variable, which is of type object.Object.
	Value object.Object
}
