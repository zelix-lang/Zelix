/*
   The Fluent Programming Language
   -----------------------------------------------------
   This code is released under the GNU GPL v3 license.
   For more information, please visit:
   https://www.gnu.org/licenses/gpl-3.0.html
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
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
