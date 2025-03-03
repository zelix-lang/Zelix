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

package property

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/property"
	queue2 "fluent/analyzer/queue"
	"fluent/ast"
	"fluent/filecode"
	"fmt"
)

// ProcessPropIdentifier processes a property identifier within a module.
// It checks for illegal access, verifies the existence of the variable,
// updates the type of the element, and handles module references.
//
// Parameters:
// - element: The expected pair element to be processed.
// - trace: The file code trace containing module information.
// - child: The AST node representing the property identifier.
//
// Returns:
// - An error3.Error indicating the result of the processing.
func ProcessPropIdentifier(
	element *queue2.ExpectedPair,
	trace *filecode.FileCode,
	child *ast.AST,
) error3.Error {
	lastPropValue, err := property.EvaluateLastPropValue(element)

	if err.Code != error3.Nothing {
		return err
	}

	// Check for illegal access
	if lastPropValue.Path != *child.File {
		return error3.Error{
			Code:   error3.IllegalPropAccess,
			Line:   element.Tree.Line,
			Column: element.Tree.Column,
		}
	}

	// Check if the variable exists within the module
	value, found := lastPropValue.Declarations[*child.Value]

	// Return the error if it is not found
	if !found {
		fmt.Println(*lastPropValue, child.Marshal(0))
		return error3.Error{
			Code:       error3.UndefinedReference,
			Additional: []string{*child.Value},
			Line:       element.Tree.Line,
			Column:     element.Tree.Column,
		}
	}

	// Check for constant reassignments
	if element.IsPropReassignment && value.IsConstant {
		return error3.Error{
			Code:   error3.ConstantReassignment,
			Line:   element.Tree.Line,
			Column: element.Tree.Column,
		}
	}

	// Change the got type to the type of the value
	oldPointers := element.Got.Type.PointerCount
	oldArrays := element.Got.Type.ArrayCount
	element.Got.Type = value.Type
	element.Got.Type.PointerCount += oldPointers
	element.Got.Type.ArrayCount += oldArrays
	element.ActualPointers += value.Type.PointerCount

	if !element.Got.Type.IsPrimitive {
		*element.LastPropValue = trace.Modules[element.Got.Type.BaseType]
	}

	// Check for modules
	if !value.Type.IsPrimitive {
		// Get the module
		mod, found := trace.Modules[value.Type.BaseType]

		// Return the error if it is not found
		if !found {
			return error3.Error{
				Code:       error3.UndefinedReference,
				Additional: []string{value.Type.BaseType},
				Line:       element.Tree.Line,
				Column:     element.Tree.Column,
			}
		}

		// Update the result
		element.Got.Value = mod
	}

	return error3.Error{}
}
