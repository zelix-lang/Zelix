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

import "fluent/filecode/types/wrapper"

// IRVariable represents a variable in the intermediate representation (IR).
type IRVariable struct {
	// Addr is the address of the variable.
	Addr string
	// Type is the type of the variable, wrapped in a TypeWrapper.
	Type *wrapper.TypeWrapper
}
