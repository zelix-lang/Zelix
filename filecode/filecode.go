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

package filecode

import (
	"fluent/filecode/function"
	"fluent/filecode/module"
)

// FileCode represents the structure of a file in the Fluent programming language.
type FileCode struct {
	// Path is the file path.
	Path string
	// Functions is a map of function names to their corresponding Function objects.
	Functions map[string]function.Function
	// Modules is a map of module names to their corresponding Module objects.
	Modules map[string]module.Module
	// Imports is a list of imported packages.
	Imports []string
	// The contents of the file that this FileCode represents.
	Contents string
}
