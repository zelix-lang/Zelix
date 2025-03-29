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
	error2 "fluent/analyzer/error"
	"fluent/filecode/function"
)

func AnalyzeMainFunction(fn *function.Function) *error2.Error {
	// Check if the function is nil
	if fn == nil {
		return &error2.Error{
			Line:   1,
			Column: 1,
			Code:   error2.NoMainFunction,
		}
	}

	// Make sure the return is nothing
	if fn.ReturnType.BaseType != "nothing" {
		return &error2.Error{
			Line:   fn.ReturnType.Trace.Line,
			Column: fn.ReturnType.Trace.Column,
			Code:   error2.MainFunctionHasReturn,
		}
	}

	// Make sure the function has no parameters
	if len(fn.Params) > 0 {
		return &error2.Error{
			Line:   fn.Params[0].Trace.Line,
			Column: fn.Params[0].Trace.Column,
			Code:   error2.MainFunctionHasParameters,
		}
	}

	// Make sure the function does not have generics
	if len(fn.Templates) > 0 {
		return &error2.Error{
			Line:   fn.Trace.Line,
			Column: fn.Trace.Column,
			Code:   error2.MainFunctionHasGenerics,
		}
	}

	return nil
}
