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

package function

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/format"
	"fluent/analyzer/object"
	"fluent/analyzer/pool"
	"fluent/analyzer/queue"
	"fluent/analyzer/rule/conditional"
	"fluent/analyzer/rule/declaration"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/rule/loop"
	"fluent/analyzer/rule/reassignment"
	"fluent/analyzer/rule/ret"
	"fluent/analyzer/rule/value"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/filecode/types/wrapper"
)

// destroyScope destroys the specified scopes and adds warnings for unused variables.
// Parameters:
// - scope: The scoped stack containing the scopes to be destroyed.
// - scopeIds: A slice of scope IDs to be destroyed.
// - warnings: The pool to collect warnings about unused variables.
// - mainScopeId: The ID of the main scope.
// - forceDeleteMainScope: A flag indicating whether to forcefully delete the main scope.
func destroyScope(
	scope *stack.ScopedStack,
	scopeIds []int,
	warnings *pool.ErrorPool,
	mainScopeId int,
	forceDeleteMainScope bool,
) {
	for _, scopeId := range scopeIds {
		if scopeId != mainScopeId && !forceDeleteMainScope {
			continue
		}

		unusedVariables := scope.DestroyScope(scopeId)

		// Add unused variable warnings
		for name, variable2 := range unusedVariables {
			warnings.AddError(&error3.Error{
				Code:       error3.UnusedVariable,
				Line:       variable2.Trace.Line,
				Column:     variable2.Trace.Column,
				Additional: []string{name},
			})
		}
	}
}

// AnalyzeFunction analyzes a given function and returns pools of errors and warnings.
// Parameters:
// - fun: The function to be analyzed.
// - trace: The file code trace information.
// - parentName: The name of the module that contains this function.
// - generics: The list of generics that the function has.
// - collectAssignments: Whether collect in a map all the variables reassigned throughout the function.
//
// Returns:
// - A pool of errors found during the analysis.
// - A pool of warnings found during the analysis.
// - A slice of AST nodes representing the assignments found in the function.
func AnalyzeFunction(
	fun function.Function,
	trace *filecode.FileCode,
	parentName string,
	generics *map[string]bool,
	scope stack.ScopedStack,
	collectAssignments bool,
	isMod bool,
) (*pool.ErrorPool, *pool.ErrorPool, *[]*ast.AST) {
	errors := pool.NewErrorPool()
	warnings := pool.NewErrorPool()
	var collectedAssignments *[]*ast.AST

	// Only initialize if collectAssignments is true
	if collectAssignments {
		collectedAssignments = &[]*ast.AST{}
	}

	// Check that the case matches snake_case
	if fun.Name != parentName && !format.CheckCase(&fun.Name, format.SnakeCase) {
		warnings.AddError(&error3.Error{
			Code:       error3.NameShouldBeSnakeCase,
			Line:       fun.Trace.Line,
			Column:     fun.Trace.Column,
			Additional: []string{fun.Name},
		})
	}

	// Check for undefined references in the return type
	err := value.AnalyzeUndefinedReference(trace, fun.ReturnType, generics)

	// Push the error to the list if necessary
	errors.AddError(err)

	// Used to determine if the function has returned a value
	hasReturned := false

	// Create a new scope for the function
	mainScopeId := scope.NewScope()
	returnType := fun.ReturnType

	// Analyze and add all parameters to the scope
	for _, param := range fun.Params {
		err, warn := AnalyzeParameter(&param.Name, &param, trace, generics)

		// Push the error to the list if necessary
		errors.AddError(err)
		warnings.AddError(warn)

		// Create the value
		val := object.Object{
			Type:   param.Type,
			Value:  nil,
			IsHeap: param.Type.PointerCount > 0,
		}

		// See if the parameter's type is a module
		if !param.Type.IsPrimitive {
			// Get the module
			module := trace.Modules[param.Type.BaseType]

			// Set the value
			val.Value = module
		}

		// Add the parameter to the scope
		scope.Append(param.Name, variable.Variable{
			Constant: true,
			Value:    val,
			Trace:    param.Trace,
		}, 0)
	}

	// Use a queue to analyze the function's block and nested blocks within it
	blockQueue := make([]queue.BlockQueueElement, 0)
	blockQueue = append(blockQueue, queue.BlockQueueElement{
		Block: &fun.Body,
		ID:    []int{mainScopeId},
	})

	// Add a main scope ID to the first element of the queue
	if isMod {
		blockQueue[0].ID = append(blockQueue[0].ID, 0) // Add the main scope ID for modules
	}

	for len(blockQueue) > 0 {
		// Get the first element in the queue
		element := blockQueue[0]
		block := element.Block
		scopeIds := element.ID
		inLoop := element.InLoop
		blockQueue = blockQueue[1:]

		// Special case: For loops (creates the block at the end)
		dontDeleteStack := false

		for _, statement := range *block.Children {
			rule := statement.Rule

			switch rule {
			case ast.Return:
				hasReturned = true
				err := ret.AnalyzeReturn(statement, trace, &scope, &returnType, element.ID)

				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.Block:
				dontDeleteStack = true
				// Create a new scope
				newScopeId := scope.NewScope()
				scopeIds = append(scopeIds, newScopeId)

				// Add the block to the queue
				blockQueue = append(blockQueue, queue.BlockQueueElement{
					Block:  statement,
					ID:     scopeIds,
					InLoop: inLoop,
				})
			case ast.Declaration:
				err, warning := declaration.AnalyzeDeclaration(statement, &scope, trace, generics, parentName, element.ID)

				// Push the error to the list if necessary
				errors.AddError(err)
				warnings.AddError(warning)
			case ast.Assignment:

				err := reassignment.AnalyzeReassignment(statement, &scope, trace, collectedAssignments, element.ID)
				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.If:
				dontDeleteStack = true
				err := conditional.AnalyzeIf(
					statement,
					trace,
					&scope,
					&blockQueue,
					scopeIds,
					inLoop,
				)
				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.While:
				dontDeleteStack = true
				err := conditional.ProcessSingleConditional(
					*statement.Children,
					trace,
					&scope,
					&blockQueue,
					scopeIds,
					true,
				)
				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.For:
				dontDeleteStack = true
				err := loop.AnalyzeFor(
					statement,
					trace,
					&scope,
					&blockQueue,
					scopeIds,
				)

				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.Break, ast.Continue:
				if !inLoop {
					errors.AddError(&error3.Error{
						Code:   error3.InvalidLoopInstruction,
						Line:   statement.Line,
						Column: statement.Column,
					})
				}
			default:
				_, err := expression.AnalyzeExpression(
					statement,
					trace,
					&scope,
					false,
					&wrapper.TypeWrapper{
						Children: &[]*wrapper.TypeWrapper{},
					},
					false,
					false,
					element.ID,
				)

				// Push the error to the list if necessary
				errors.AddError(err)
			}
		}

		if !dontDeleteStack && len(scopeIds) > 0 && scopeIds[0] != 0 {
			// Avoid destroying the main scope
			destroyScope(&scope, scopeIds, warnings, mainScopeId, false)
		}
	}

	// Destroy the main scope at the end
	destroyScope(&scope, []int{mainScopeId}, warnings, mainScopeId, true)

	// Make sure that the function has returned a value
	if !fun.IsStd && fun.Name != "heap_alloc" && fun.ReturnType.BaseType != "nothing" && !hasReturned {
		errors.AddError(&error3.Error{
			Code:   error3.MustReturnAValue,
			Line:   fun.Trace.Line,
			Column: fun.Trace.Column,
		})
		return errors, warnings, collectedAssignments
	}

	return errors, warnings, collectedAssignments
}
