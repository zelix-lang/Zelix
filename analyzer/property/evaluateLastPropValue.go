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
	"fluent/analyzer/queue"
	"fluent/filecode/module"
)

// EvaluateLastPropValue evaluates the last property value of the given element.
// If the element is a property access and the last property value is not nil,
// it attempts to cast the value to a module.Module. If the cast is successful,
// it returns the module. Otherwise, it returns an error indicating invalid property access.
//
// Parameters:
//   - element: A pointer to a queue.ExpectedPair which contains the property access information.
//
// Returns:
//   - A pointer to a module.Module if the cast is successful.
//   - An error3.Error indicating invalid property access if the cast fails or the last property value is nil.
func EvaluateLastPropValue(element *queue.ExpectedPair) (*module.Module, error3.Error) {
	if element.IsPropAccess {
		if element.LastPropValue == nil {
			return nil, error3.Error{
				Code:   error3.InvalidPropAccess,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		var convert interface{}
		convert = *element.LastPropValue

		// Cast the last property value to a module
		mod, castOk := convert.(module.Module)

		if !castOk {
			nMod, nCastOk := convert.(*module.Module)

			if nCastOk {
				castOk = true
				mod = *nMod
			}
		}

		if !castOk {
			return nil, error3.Error{
				Code:   error3.InvalidPropAccess,
				Line:   element.Tree.Line,
				Column: element.Tree.Column,
			}
		}

		return &mod, error3.Error{}
	}

	return nil, error3.Error{
		Code:   error3.InvalidPropAccess,
		Line:   element.Tree.Line,
		Column: element.Tree.Column,
	}
}
