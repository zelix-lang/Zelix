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

package pkg

// Package represents a Fluent project package with metadata.
type Package struct {
	// Name is the name of the package.
	Name string
	// Version is the version of the package.
	Version string
	// Description is a brief description of the package.
	Description string
	// Author is the author of the package.
	Author string
	// License is the license under which the package is distributed.
	License string
	// Entry is the entry point of the package.
	Entry string
}
