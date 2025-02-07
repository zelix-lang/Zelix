/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package analyzer

import (
	error3 "fluent/analyzer/error"
	"fluent/analyzer/rule"
	"fluent/filecode"
	"fluent/filecode/function"
	"fluent/filecode/module"
	trace2 "fluent/filecode/trace"
	"fluent/logger"
	error2 "fluent/message/error"
	"fluent/message/warn"
	"fluent/state"
	"fluent/util"
	"fmt"
	"os"
	"strings"
)

// checkImportRedefinition checks for redefinition of imported functions or modules.
// If a redefinition is found, it logs the details and exits the program.
//
// Parameters:
// - collection: A map of existing functions or modules.
// - name: The name of the function or module being checked.
// - value: The function or module being checked.
// - file: The file where the function or module is defined.
// - entry: A map where the key is the file path and the value is the FileCode object.
func checkImportRedefinition[T function.Function | module.Module](
	collection map[string]T,
	name string,
	value T,
	file *filecode.FileCode,
	entry *map[string]filecode.FileCode,
) {
	public := false
	var valueTrace trace2.Trace
	entryDeref := *entry

	// Extract fields from value using type switch
	switch v := any(value).(type) {
	case function.Function:
		public = v.Public
		valueTrace = v.Trace
	case module.Module:
		public = v.Public
		valueTrace = v.Trace
	}

	if definedVal, ok := collection[name]; ok && public {
		var trace trace2.Trace
		var path string

		// Extract fields from definedVal using a type switch
		switch v := any(definedVal).(type) {
		case function.Function:
			trace = v.Trace
			path = v.Path
		case module.Module:
			trace = v.Trace
			path = v.Path
		}

		error2.Redefinition(name)
		fmt.Println(
			util.BuildDetails(
				&file.Contents,
				&file.Path,
				valueTrace.Line,
				valueTrace.Column,
				true,
			),
		)

		logger.Info("'" + name + "' was previously defined here:")

		originalContents := entryDeref[path].Contents
		originalPath := entryDeref[path].Path
		fmt.Println(
			util.BuildDetails(
				&originalContents,
				&originalPath,
				trace.Line,
				trace.Column,
				true,
			),
		)

		os.Exit(1)
	}
}

// AnalyzeCode analyzes the provided code files based on their dependencies.
// It uses a priority-based queue to analyze files without dependencies first.
//
// Parameters:
// - entry: A map where the key is the file path and the value is the FileCode object.
// - mainPath: The path to the main file that should be analyzed last.
// - silent: A boolean that indicates if the analysis should be silent.
func AnalyzeCode(entry map[string]filecode.FileCode, mainPath string, silent bool) {
	// Use a priority-based queue to analyze the files that do not have
	// dependencies first
	var queue []filecode.FileCode

	// A map that stores the seen paths (a map for fast lookup)
	seen := make(map[string]bool)

	// Use a breadth-first push queue to push the remaining files
	pushQueue := make([]filecode.FileCode, 0)

	// Iterate over all the entries to prioritize the files
	for path, file := range entry {
		if path == mainPath {
			// The main path is analyzed at the end
			continue
		}

		// Check if the file has dependencies
		if len(file.Imports) > 0 {
			pushQueue = append(pushQueue, file)
			continue
		}

		seen[path] = true
		// Push the file to the end of the queue
		queue = append(queue, file)
	}

	// Push all the remaining files to the queue
	for len(pushQueue) > 0 {
		// Pop the first element
		file := pushQueue[0]
		pushQueue = pushQueue[1:]

		// Check if the file has dependencies
		hasDependencies := false
		for _, importPath := range file.Imports {
			if !seen[importPath] {
				hasDependencies = true
				break
			}
		}

		if hasDependencies {
			// Push the file to the push queue
			pushQueue = append(pushQueue, file)
		} else {
			// Push the file to the end of the queue
			queue = append(queue, file)
			seen[file.Path] = true
		}
	}

	// Append the main file to the end of the queue
	queue = append(queue, entry[mainPath])

	// Analyze the files
	for _, file := range queue {
		if !silent {
			state.Emit(state.Analyzing, util.FileName(&file.Path))
		}

		// Append all the imported modules and functions to the file
		for _, importPath := range file.Imports {
			importedFile := entry[importPath]

			// Append the imported functions
			for _, fun := range importedFile.Functions {
				// Check for redefinitions
				checkImportRedefinition(
					file.Functions,
					fun.Name,
					fun,
					&file,
					&entry,
				)

				file.Functions[fun.Name] = fun
			}

			// Append the imported modules
			for _, mod := range importedFile.Modules {
				// Check for redefinitions
				checkImportRedefinition(
					file.Modules,
					mod.Name,
					mod,
					&file,
					&entry,
				)

				file.Modules[mod.Name] = mod
			}
		}

		errors, warnings := rule.AnalyzeFileCode(file)

		if errors.Count > 0 {
			state.FailAllSpinners()
		} else if warnings.Count > 0 {
			state.WarnAllSpinners()
		}

		// Use a strings.Builder to build the error and warning messages
		var errorMessage strings.Builder
		var warningMessage strings.Builder

		// Print the errors
		for _, err := range errors.Errors {
			switch err.Code {
			case error3.ParamTypeNothing:
				errorMessage.WriteString(error2.NothingParam())
			case error3.Redefinition:
				errorMessage.WriteString(error2.Redefinition(err.Additional[0]))
			case error3.TypeMismatch:
				errorMessage.WriteString(error2.TypeMismatch(err.Additional[0], err.Additional[1]))
			case error3.UndefinedReference:
				errorMessage.WriteString(error2.UndefinedReference(err.Additional[0]))
			case error3.InvalidDereference:
				errorMessage.WriteString(error2.InvalidDereference())
			case error3.MustReturnAValue:
				errorMessage.WriteString(error2.MustReturnValue())
			case error3.DataOutlivesStack:
				errorMessage.WriteString(error2.DataOutlivesStack())
			case error3.ParameterCountMismatch:
				errorMessage.WriteString(error2.ParamCountMismatch(err.Additional[0]))
			case error3.CannotInferType:
				errorMessage.WriteString(error2.CannotInferType())
			case error3.ShouldNotReturn:
				errorMessage.WriteString(error2.ShouldNotReturn())
			case error3.CannotTakeAddress:
				errorMessage.WriteString(error2.CannotTakeAddress())
			case error3.InvalidPropAccess:
				errorMessage.WriteString(error2.InvalidPropAccess())
			case error3.IllegalPropAccess:
				errorMessage.WriteString(error2.IllegalPropAccess())
			case error3.ConstantReassignment:
				errorMessage.WriteString(error2.ConstantReassignment())
			default:
			}

			// Write the details
			errorMessage.WriteString(util.BuildDetails(&file.Contents, &file.Path, err.Line, err.Column, true))
		}

		for _, warning := range warnings.Errors {
			switch warning.Code {
			case error3.NameShouldBeSnakeCase:
				warningMessage.WriteString(warn.SnakeCase(warning.Additional[0]))
			case error3.UnusedVariable:
				warningMessage.WriteString(warn.UnusedVariable(warning.Additional[0]))
			default:
			}

			// Write the details
			warningMessage.WriteString(util.BuildDetails(&file.Contents, &file.Path, warning.Line, warning.Column, false))
		}

		if errors.Count > 0 {
			fmt.Print(errorMessage.String())
			os.Exit(1)
		}

		if warnings.Count > 0 {
			fmt.Print(warningMessage.String())
		}

		// Remove foreign functions
		for i, fun := range file.Functions {
			if fun.Path != file.Path {
				delete(file.Functions, i)
			}
		}

		// Remove foreign modules
		for i, mod := range file.Modules {
			if mod.Path != file.Path {
				delete(file.Modules, i)
			}
		}

		state.PassAllSpinners()
	}
}
