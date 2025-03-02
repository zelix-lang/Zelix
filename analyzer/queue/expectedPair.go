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

import (
	"fluent/analyzer/object"
	"fluent/ast"
	"fluent/filecode/types/wrapper"
)

// ExpectedPair represents a pair of expected and actual values
// along with the corresponding AST node.
type ExpectedPair struct {
	Expected           *wrapper.TypeWrapper // The expected type
	Got                *object.Object       // The actual object
	Tree               *ast.AST             // The AST node
	HasMetDereference  bool                 // Whether a dereference token has been met
	ActualPointers     int                  // The number of pointers in the actual object
	HeapRequired       bool                 // Whether this value must be heap-allocated
	ModRequired        bool                 // Whether this value must evaluate to a module
	IsPropAccess       bool                 // Whether this value is a property access
	IsArithmetic       bool                 // Whether this value must evaluate a decimal or integer
	IsPropReassignment bool                 // Whether this value is a property reassignment
	LastPropValue      *interface{}         // The last property value
	IsParam            bool                 // Whether this value is a parameter
}
