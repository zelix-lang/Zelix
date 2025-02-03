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
	"fluent/logger"
	error2 "fluent/message/error"
	"fluent/state"
	"fluent/util"
	"os"
)

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
			if !silent {
				state.Emit(state.Processing, util.FileName(&file.Path))
			}

			// Push the file to the end of the queue
			queue = append(queue, file)
			seen[file.Path] = true

			state.PassAllSpinners()
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
			for _, function := range importedFile.Functions {
				// Check for redefinitions
				if definedFun, ok := file.Functions[function.Name]; ok && function.Public {
					error2.Redefinition(definedFun.Name)
					util.BuildAndPrintDetails(
						&file.Contents,
						&file.Path,
						definedFun.Trace.Line,
						definedFun.Trace.Column,
						true,
					)

					logger.Info("'" + function.Name + "' was previously defined here:")

					originalContents := entry[function.Path].Contents
					originalPath := entry[function.Path].Path
					util.BuildAndPrintDetails(
						&originalContents,
						&originalPath,
						function.Trace.Line,
						function.Trace.Column,
						true,
					)

					os.Exit(1)
				}

				file.Functions[function.Name] = function
			}

			// Append the imported modules
			for _, module := range importedFile.Modules {
				// Check for redefinitions
				if definedMod, ok := file.Modules[module.Name]; ok && module.Public {
					error2.Redefinition(definedMod.Name)
					util.BuildAndPrintDetails(
						&file.Contents,
						&file.Path,
						definedMod.Trace.Line,
						definedMod.Trace.Column,
						true,
					)

					logger.Info("'" + module.Name + "' was previously defined here:")

					originalContents := entry[module.Path].Contents
					originalPath := entry[module.Path].Path

					util.BuildAndPrintDetails(
						&originalContents,
						&originalPath,
						module.Trace.Line,
						module.Trace.Column,
						true,
					)

					os.Exit(1)
				}

				file.Modules[module.Name] = module
			}
		}

		errors := rule.AnalyzeFileCode(file)
		if errors.Count > 0 {
			state.FailAllSpinners()
		}

		// Print the errors
		for _, err := range errors.Errors {
			switch err.Code {
			case error3.Redefinition:
				error2.Redefinition(err.Additional)
			case error3.TypeMismatch:
				error2.TypeMismatch()
			case error3.UndefinedReference:
				error2.UndefinedReference(err.Additional)
			case error3.InvalidDereference:
				error2.InvalidDereference()
			case error3.MustReturnAValue:
				error2.MustReturnValue()
			case error3.DataOutlivesStack:
				error2.DataOutlivesStack()
			case error3.ParameterCountMismatch:
				error2.ParamCountMismatch(err.Additional)
			case error3.CannotInferType:
				error2.CannotInferType()
			default:
			}

			// Print the details
			util.BuildAndPrintDetails(&file.Contents, &file.Path, err.Line, err.Column, true)
		}

		if errors.Count > 0 {
			os.Exit(1)
		}

		// Remove foreign functions
		for i, function := range file.Functions {
			if function.Path != file.Path {
				delete(file.Functions, i)
			}
		}

		// Remove foreign modules
		for i, module := range file.Modules {
			if module.Path != file.Path {
				delete(file.Modules, i)
			}
		}

		state.PassAllSpinners()
	}
}
