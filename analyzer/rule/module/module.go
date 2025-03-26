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

package module

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/object"
	"fluent/analyzer/pool"
	"fluent/analyzer/rule/expression"
	function2 "fluent/analyzer/rule/function"
	"fluent/analyzer/rule/value"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/filecode"
	"fluent/filecode/module"
	"fluent/filecode/types/wrapper"
	"fluent/logger"
	"fmt"
	"strings"
)

// AnalyzeModule analyzes a given module and its functions, checking for errors and warnings.
// It returns two error pools: one for errors and one for warnings.
//
// Parameters:
// - mod: The module to be analyzed.
// - trace: The file code trace for the module.
//
// Returns:
// - *pool.ErrorPool: A pool containing all the errors found during the analysis.
// - *pool.ErrorPool: A pool containing all the warnings found during the analysis.
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
			Type: wrapper.TypeWrapper{
				BaseType: mod.Name,
				Children: &[]*wrapper.TypeWrapper{},
			},
		},
	})

	// Analyze all functions in the module
	for _, fun := range mod.Functions {
		// Make sure that the function does not have generics
		if len(fun.Templates) != 0 {
			globalErrors.AddError(&error3.Error{
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
			*fun,
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
					globalErrors.AddError(&error3.Error{
						Line:   declaration.Trace.Line,
						Column: declaration.Trace.Column,
						Code:   error3.ValueNotAssigned,
					})
				}
			}
		}
	}

	// Use a queue to analyze the module's declarations
	// to find circular module dependencies
	queue := make([]module.Declaration, 0)
	// Save the seen modules in a slice for printing
	seenModules := []string{mod.Name}
	// Save the seen modules in a map for fast lookup
	seenModulesMap := map[string]bool{
		mod.Name: true,
	}

	// Used to determine if a value that has been not initialized has already been caught
	initializedCaught := false

	// Check if the module has incomplete declarations
	for _, declaration := range mod.Declarations {
		// Check for undefined references
		err := value.AnalyzeUndefinedReference(trace, declaration.Type, &mod.Templates)
		queue = append(queue, declaration)

		if err != nil {
			globalErrors.AddError(err)
		}

		if !initializedCaught && declaration.IsIncomplete {
			initializedCaught = true
			// Check that the module has a constructor
			if _, ok := mod.Functions[mod.Name]; !ok {
				globalErrors.AddError(&error3.Error{
					Line:   declaration.Trace.Line,
					Column: declaration.Trace.Column,
					Code:   error3.ValueNotAssigned,
				})
			}
		}

		// Analyze the expression's value
		_, err = expression.AnalyzeExpression(
			declaration.Value,
			trace,
			&variables,
			false,
			&declaration.Type,
			false,
			true,
		)

		if err != nil {
			globalErrors.AddError(err)
		}
	}

	// Find circular dependencies
	for len(queue) > 0 {
		// Get the first element of the queue
		element := queue[0]
		queue = queue[1:]

		// Check if the type is a module
		if element.Type.IsPrimitive {
			continue
		}

		// Get the base type of the declaration
		baseType := element.Type.BaseType

		// Check if the module has already been seen
		if seenModulesMap[baseType] {
			// Use a strings.builder to build the full import chain in an efficient manner
			builder := strings.Builder{}
			spaces := 0

			for _, modName := range seenModules {
				builder.WriteString(
					logger.BuildInfo(
						fmt.Sprintf(
							"%s-> %s",
							strings.Repeat("  ", spaces),
							modName,
						),
					),
				)
				spaces++
			}

			// Also write the current module
			builder.WriteString(
				logger.BuildInfo(
					fmt.Sprintf(
						"%s-> %s",
						strings.Repeat("  ", spaces),
						mod.Name,
					),
				),
			)

			globalErrors.AddError(&error3.Error{
				Line:       element.Trace.Line,
				Column:     element.Trace.Column,
				Code:       error3.CircularModuleDependency,
				Additional: []string{builder.String()},
			})

			break
		}

		// Add the module to the seen modules
		seenModules = append(seenModules, baseType)
		seenModulesMap[baseType] = true

		// Get the module
		mod, ok := trace.Modules[baseType]

		if !ok {
			// Error is caught in further analysis for that specific module
			continue
		}

		// Add all the declarations to the queue
		for _, declaration := range mod.Declarations {
			queue = append(queue, declaration)
		}
	}

	return globalErrors, globalWarnings
}
