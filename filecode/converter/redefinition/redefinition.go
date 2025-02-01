/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package redefinition

import (
	"fluent/filecode/function"
	"fluent/filecode/module"
	trace2 "fluent/filecode/trace"
	"fluent/logger"
	error2 "fluent/message/error"
	"fluent/util"
	"os"
)

func CheckRedefinition[T function.Function | module.Module](
	definedValues map[string]T,
	name string,
	entity any,
	contents string,
	path string,
) {
	var trace trace2.Trace

	if fn, ok := entity.(function.Function); ok {
		trace = fn.Trace
	} else if mod, ok := entity.(module.Module); ok {
		trace = mod.Trace
	} else {
		logger.Error("Unknown entity type")
		os.Exit(1)
	}

	if _, ok := definedValues[name]; ok {
		error2.Redefinition(name)
		util.BuildAndPrintDetails(
			&contents,
			&path,
			trace.Line,
			trace.Column,
			true,
		)

		logger.Info("'" + name + "' was previously defined here:")
		util.BuildAndPrintDetails(
			&contents,
			&path,
			trace.Line,
			trace.Column,
			true,
		)

		os.Exit(1)
	}
}
