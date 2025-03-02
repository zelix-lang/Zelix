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

package function

import (
	"fluent/filecode/trace"
	"fluent/filecode/types/wrapper"
)

// Param represents a function parameter with a name, type, and trace information.
type Param struct {
	// Type is the type of the parameter, wrapped in a TypeWrapper.
	Type wrapper.TypeWrapper
	// Trace contains trace information for the parameter.
	Trace trace.Trace
	// Name is the name of the parameter.
	Name string
}
