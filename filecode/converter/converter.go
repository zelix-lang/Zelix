/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package converter

import (
	"fluent/ansi"
	ast2 "fluent/ast"
	"fluent/filecode"
	"fluent/filecode/converter/function"
	module2 "fluent/filecode/converter/module"
	"fluent/filecode/converter/redefinition"
	function2 "fluent/filecode/function"
	"fluent/filecode/module"
	"fluent/lexer"
	"fluent/logger"
	"fluent/parser"
	"fluent/state"
	"fluent/util"
	"os"
	"strings"
)

// The standard library path
var stdPath = os.Getenv("FLUENT_STD_PATH")

// The system's path separator
const pathSeparator = string(os.PathSeparator)

type queueElement struct {
	path     *string
	trace    *ast2.AST
	contents *string
}

// ConvertToFileCode converts the given entry file and all its imports to a map of FileCode.
// It checks for the existence of the standard library path, reads files, tokenizes them, parses them into ASTs,
// and then converts the ASTs into FileCode structures. It also handles circular import detection.
//
// Parameters:
// - entry: The entry file path to start the conversion.
//
// Returns:
// - A map where the keys are file paths and the values are FileCode structures.
func ConvertToFileCode(entry string) map[string]filecode.FileCode {
	// Check that the stdlib path exists
	if stdPath == "" || !util.DirExists(stdPath) {
		logger.Error("The FLUENT_STD_PATH environment variable is not set")
		logger.Help("Try reinstalling the CLI")
		os.Exit(1)
	}

	// Use a queue to convert the file and all of its imports
	queue := []queueElement{
		{
			path: &entry,
		},
	}

	// Save the seen imports in a map (for O(1) lookup) to detect circular imports
	seenImports := map[string]bool{}
	result := make(map[string]filecode.FileCode)

	// Save the already-processed std imports
	processedStdImports := map[string]bool{}

	for len(queue) > 0 {
		// Get the first element of the queue
		element := queue[0]
		queue = queue[1:]

		path := element.path
		trace := element.trace
		elContents := element.contents

		isStd := strings.HasPrefix(*path, stdPath)

		// Detect circular imports
		if seenImports[*path] {
			logger.Error("Circular import detected")
			logger.Info("Full import chain:")

			spaces := 0
			for importPath := range seenImports {
				if *path == importPath {
					logger.Info(
						ansi.Colorize(
							ansi.BoldBrightRed,
							strings.Repeat("  ", spaces)+"-> "+util.DiscardCwd(&importPath)+" (Circular)",
						),
					)
				} else {
					logger.Info(strings.Repeat("  ", spaces) + "-> " + util.DiscardCwd(&importPath))
				}

				spaces++
			}

			// Also print the current circular import's details
			logger.Info(
				ansi.Colorize(
					ansi.BoldBrightRed,
					strings.Repeat("  ", spaces)+"-> "+util.DiscardCwd(path)+" (Circular)",
				),
			)

			logger.Info("Full details:")
			util.BuildAndPrintDetails(elContents, path, trace.Line, trace.Column, true)

			os.Exit(1)
		}

		if !isStd {
			seenImports[*path] = true
		}

		// Read the file
		contents := util.ReadFile(*path)

		// Lex the file
		state.Emit(state.Lexing, util.FileName(path))
		tokens, lexerError := lexer.Lex(contents, *path)

		if lexerError.Message != "" {
			state.FailAllSpinners()
			// Build and print the error
			util.PrintError(&contents, path, &lexerError.Message, lexerError.Line, lexerError.Column)
			os.Exit(1)
		}

		state.PassAllSpinners()
		// Parse the tokens to an AST
		ast, parsingError := parser.Parse(tokens, *path)

		if parsingError.IsError() {
			state.FailAllSpinners()
			errorMessage := util.BuildMessageFromParsingError(parsingError)

			// Build and print the error
			util.PrintError(&contents, path, &errorMessage, parsingError.Line, parsingError.Column)
			os.Exit(1)
		}

		state.PassAllSpinners()
		state.Emit(state.Parsing, util.FileName(path))
		code := filecode.FileCode{
			Path:      *path,
			Functions: make(map[string]function2.Function),
			Modules:   make(map[string]module.Module),
			Imports:   make([]string, 0),
			Contents:  contents,
		}

		// Traverse the AST to convert it to a FileCode
		for _, child := range *ast.Children {
			rule := child.Rule

			switch rule {
			case ast2.Import:
				// Get the path
				importPath := *(*child.Children)[0].Value

				isStd := strings.HasPrefix(importPath, "@std")

				if isStd {
					if processedStdImports[importPath] {
						continue
					}

					processedStdImports[importPath] = true
					importPath = strings.Replace(importPath, "@std", stdPath, 1)
					importPath = strings.ReplaceAll(importPath, "::", pathSeparator)
				} else {
					// Get the file's directory
					dir := util.GetDir(*path)

					// Join the directory with the path
					importPath = dir + pathSeparator + importPath
				}

				if !strings.HasSuffix(importPath, ".fluent") {
					importPath += ".fluent"
				}

				// Append the path to the code's imports
				code.Imports = append(code.Imports, importPath)

				// Queue the path
				queue = append(queue, queueElement{
					path:     &importPath,
					contents: &contents,
					trace:    child,
				})
			case ast2.Function:
				// Convert to a Function wrapper
				fn := function.ConvertFunction(child)

				// Check for redefinitions
				redefinition.CheckRedefinition(code.Functions, fn.Name, fn, contents, *path)

				code.Functions[fn.Name] = fn
			case ast2.Module:
				// Convert to a Function wrapper
				mod := module2.ConvertModule(child, contents)

				// Check for redefinitions
				redefinition.CheckRedefinition(code.Modules, mod.Name, mod, contents, *path)

				code.Modules[mod.Name] = mod
			default:
			}
		}

		state.PassAllSpinners()
		result[*path] = code
	}

	return result
}
