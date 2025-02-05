/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package queue

import (
	"fluent/analyzer/object"
	"fluent/ast"
	"fluent/filecode/types"
)

// ExpectedPair represents a pair of expected and actual values
// along with the corresponding AST node.
type ExpectedPair struct {
	Expected          *types.TypeWrapper // The expected type
	Got               *object.Object     // The actual object
	Tree              *ast.AST           // The AST node
	HasMetDereference bool               // Whether a dereference token has been met
	ActualPointers    int                // The number of pointers in the actual object
	HeapRequired      bool               // Whether this value must be heap-allocated
	ModRequired       bool               // Whether this value must evaluate to a module
	IsPropAccess      bool               // Whether this value is a property access
}
