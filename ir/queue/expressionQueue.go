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

package queue

import "fluent/ast"

// PendingIRMarshal represents a structure for pending IR marshaling.
type PendingIRMarshal struct {
	// Input is the AST input for the pending IR marshal.
	Input *ast.AST
	// IsParam indicates if the input is a parameter.
	IsParam bool
	// HasProcessedParams indicates if the parameters have been processed.
	HasProcessedParams bool
	// Counter is used to count the number of operations.
	Counter int
}
