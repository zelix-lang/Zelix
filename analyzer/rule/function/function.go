/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package function

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/format"
	"fluent/analyzer/object"
	"fluent/analyzer/pool"
	"fluent/analyzer/rule/expression"
	"fluent/analyzer/rule/ret"
	"fluent/analyzer/rule/value"
	"fluent/analyzer/stack"
	"fluent/analyzer/variable"
	"fluent/ast"
	"fluent/filecode"
	"fluent/filecode/function"
)

func AnalyzeFunction(fun function.Function, trace *filecode.FileCode) (*pool.ErrorPool, *pool.ErrorPool) {
	errors := pool.NewErrorPool()
	warnings := pool.NewErrorPool()

	// Check that the case matches snake_case
	if !format.CheckCase(&fun.Name, format.SnakeCase) {
		warnings.AddError(error3.Error{
			Code:       error3.NameShouldBeSnakeCase,
			Line:       fun.Trace.Line,
			Column:     fun.Trace.Column,
			Additional: []string{fun.Name},
		})
	}

	// Check for undefined references in the return type
	err := value.AnalyzeUndefinedReference(trace, fun.ReturnType, &fun.Templates)

	// Push the error to the list if necessary
	errors.AddError(err)

	// Used to determine if the function has returned a value
	hasReturned := false

	// Create a new scope for the function
	scope := stack.ScopedStack{}
	scope.NewScope()

	returnType := fun.ReturnType

	// Analyze and add all parameters to the scope
	for name, param := range fun.Params {
		err, warn := AnalyzeParameter(&name, &param, trace, &fun.Templates)

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
		scope.Append(name, variable.Variable{
			Constant: true,
			Value:    val,
		})
	}

	// Use a queue to analyze the function's block and nested blocks within it
	blockQueue := make([]ast.AST, 0)
	blockQueue = append(blockQueue, fun.Body)

	for len(blockQueue) > 0 {
		// Get the first element in the queue
		block := blockQueue[0]
		blockQueue = blockQueue[1:]

		for _, statement := range *block.Children {
			rule := statement.Rule

			switch rule {
			case ast.Return:
				hasReturned = true
				err := ret.AnalyzeReturn(statement, trace, &scope, &returnType)

				// Push the error to the list if necessary
				errors.AddError(err)
			case ast.Block:
				// Create a new scope
				scope.NewScope()
				// Add the block to the queue
				blockQueue = append(blockQueue, *statement)
			case ast.Assignment:

			default:
				_, err := expression.AnalyzeExpression(statement, trace, &scope, false)

				// Push the error to the list if necessary
				errors.AddError(err)
			}
		}
	}

	// Make sure that the function has returned a value
	if fun.ReturnType.BaseType != "nothing" && !hasReturned {
		return errors, warnings
	}

	unusedVariables := scope.DestroyScope()

	// Add unused variable warnings
	for _, variable2 := range unusedVariables {
		warnings.AddError(error3.Error{
			Code:       error3.UnusedVariable,
			Line:       fun.Trace.Line,
			Column:     fun.Trace.Column,
			Additional: []string{*variable2},
		})
	}

	return errors, warnings
}
