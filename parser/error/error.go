/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent l -f` for details.
*/

package error

import "fluent/ast"

// Error represents a parsing error with details about its location and expected rules.
type Error struct {
	Line     int        // Line number where the error occurred.
	Column   int        // Column number where the error occurred.
	File     *string    // File name where the error occurred.
	Expected []ast.Rule // Expected rules at the error location.
}

// IsError checks if the Error instance represents an error by verifying if there are any expected rules.
func (error Error) IsError() bool {
	return len(error.Expected) > 0
}
