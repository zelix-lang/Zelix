/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package module

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/pool"
	function2 "fluent/analyzer/rule/function"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/filecode"
	"fluent/filecode/module"
	"fluent/filecode/types"
	"fluent/util"
)

func AnalyzeModule(
	mod module.Module,
	trace *filecode.FileCode,
) (*pool.ErrorPool, *pool.ErrorPool) {
	globalErrors := pool.NewErrorPool()
	globalWarnings := pool.NewErrorPool()

	// Create a new scope for the module's functions
	variables := stack.ScopedStack{
		Scopes: make(map[int]stack.Stack),
	}

	variables.NewScope()

	// Add "this" to the variables
	variables.Append("this", variable.Variable{
		Constant: true,
		Value: object.Object{
			Value: mod,
			Type: types.TypeWrapper{
				BaseType: mod.Name,
				Children: &[]*types.TypeWrapper{},
			},
		},
	})

	// Analyze all functions in the module
	for _, fun := range mod.Functions {
		// Make sure that the function does not have generics
		if len(fun.Templates) != 0 {
			globalErrors.AddError(error3.Error{
				Line:   fun.Trace.Line,
				Column: fun.Trace.Column,
				Code:   error3.ShouldNotHaveGenerics,
			})

			continue
		}

		// Check for constructors
		isConstructor := mod.Name == fun.Name

		// Analyze the function
		errors, warnings, reassignments := function2.AnalyzeFunction(
			fun,
			trace,
			mod.Name,
			&mod.Templates,
			variables,
			isConstructor,
		)

		globalErrors.Extend(errors.Errors)
		globalWarnings.Extend(warnings.Errors)

		if isConstructor {
			// Save the assigned props in a map
			assignedProps := make(map[string]bool)

			for _, reassignment := range *reassignments {
				// Get the prop name
				children := *reassignment.Children
				nameExpr := *children[1]
				nameChildren := *nameExpr.Children
				propName := *nameChildren[0].Value
				assignedProps[propName] = true
			}

			// Check for uninitialized props
			for name, declaration := range mod.Declarations {
				if !declaration.IsIncomplete {
					continue
				}

				if _, ok := assignedProps[name]; !ok {
					globalErrors.AddError(error3.Error{
						Line:   declaration.Trace.Line,
						Column: declaration.Trace.Column,
						Code:   error3.ValueNotAssigned,
					})
				}
			}
		}
	}

	// Check if the module has incomplete declarations
	for _, declaration := range mod.Declarations {
		if !declaration.IsIncomplete {
			continue
		}

		// Check that the module has a constructor
		if _, ok := mod.Functions[mod.Name]; !ok {
			globalErrors.AddError(error3.Error{
				Line:   mod.Trace.Line,
				Column: mod.Trace.Column,
				Code:   error3.ValueNotAssigned,
				Additional: []string{
					util.BuildDetails(
						&trace.Contents,
						&trace.Path,
						declaration.Trace.Line,
						declaration.Trace.Column,
						true,
					),
				},
			})
		}

		break
	}

	return globalErrors, globalWarnings
}
