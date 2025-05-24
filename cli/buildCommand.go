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

package cli

import (
	"fluent/ansi"
	"fluent/filecode/converter"
	"fluent/filecode/module"
	"fluent/ir"
	"fluent/ir/pool"
	"fluent/logger"
	"fluent/state"
	"fluent/util"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var includePath = fmt.Sprintf("%s/include", converter.StdPath)
var includePathPOSIX = fmt.Sprintf("%s/posix", includePath)
var includePathWin = fmt.Sprintf("%s/win", includePath)
var fluentExtensionRegex = regexp.MustCompile("\\.fluent$")
var isWindows = runtime.GOOS == "windows"
var isPOSIX = !isWindows

// BuildCommand compiles the given Fluent project into an executable
func BuildCommand(context *cli.Command, callCompiler bool) string {
	// Check if the user is on Windows
	if !isPOSIX {
		logger.Warn("fluentc does not (and might never) support Windows.")
		logger.Warn("Please use the fluent run to use the interpreter instead.")
		logger.Warn("Or use MinGW or WSL to run fluentc.")
	}

	// Get the optimization level (validation handled by backend)
	optimizationLevel := context.Int("optimization")

	fmt.Print(ansi.Colorize(ansi.BoldBrightYellow, "⚠️ Checking if fluentc is installed....\r"))

	// Invoke a system command to check if fluentc/fluenti is installed
	var cmdName string
	if callCompiler {
		cmdName = "fluentc"
	} else {
		cmdName = "fluenti"
	}

	cmd := exec.Command(cmdName, "--help")
	err := cmd.Run()
	if err != nil {
		// Print the whole message
		fmt.Print(ansi.Colorize(ansi.BoldBrightYellow, "⚠️ Checking if "))
		if callCompiler {
			fmt.Print(ansi.Colorize(ansi.BoldBrightYellow, "fluentc is installed...."))
		} else {
			fmt.Print(ansi.Colorize(ansi.BoldBrightYellow, "fluenti is installed...."))
		}

		logger.Error("The Fluent Compiler/Interpreter is not installed.")
		logger.Info(
			"Please install it by downloading the necessary",
			"binaries from the official repository.",
		)
		os.Exit(1)
	} else {
		// Remove the message
		fmt.Print("                                        \r")
	}

	fileCodes, fileCodesMap, dirPath, originalPath := CheckCommand(context)

	// Use a global builder to build the whole program into a single IR file
	globalBuilder := strings.Builder{}

	// Keep a counter of all the file codes that have been processed
	fileCodeCount := 0

	// A map of already-defined values used for tracing lines and columns
	traceCounters := pool.NumPool{
		Storage: make(map[int]string),
		Counter: make(map[int]int),
	}
	// Keep track of used strings (Saved in reserved spaces of memory)
	usedStrings := pool.StringPool{
		Storage: make(map[string]string),
		Counter: make(map[int]int),
		Prefix:  "__str__",
	}
	traceStrings := make(map[string]*string)
	// Used to store precomputed counters for functions' and
	// modules' names
	nameCounters := make(map[string]*map[string]*string)
	modulePropCounters := make(map[*module.Module]*util.OrderedMap[string, *string])
	// Save in a map the files that have an external
	// implementation to avoid recompiling them
	externalImpl := make(map[string]bool)

	// Precompute the counters for the names
	for _, fileCode := range fileCodes {
		// Check if this file has an external implementation
		if strings.HasPrefix(fileCode.Path, converter.StdPath) {
			var relativePath string

			// Check for POSIX-Compliant systems
			if isPOSIX {
				relativePath = strings.Replace(
					fileCode.Path,
					converter.StdPath,
					includePathPOSIX,
					1,
				)
			} else {
				relativePath = strings.Replace(
					fileCode.Path,
					converter.StdPath,
					includePathWin,
					1,
				)
			}

			relativePath = fluentExtensionRegex.ReplaceAllString(relativePath, ".ll")
			if util.FileExists(relativePath) {
				fileName := util.FileName(&fileCode.Path)
				externalImpl[fileCode.Path] = true

				fmt.Println(
					ansi.Colorize(
						ansi.BoldBrightYellow,
						fmt.Sprintf(
							"🔂 Skipped %s (External impl available)",
							fileName,
						),
					),
				)

				// Add the std instruction to the global builder
				globalBuilder.WriteString("link \"")
				globalBuilder.WriteString(relativePath)
				globalBuilder.WriteString("\"\n")
				continue
			}
		}

		// Make sure the map is initialized
		nameCounter, ok := nameCounters[fileCode.Path]

		if !ok {
			nameCounter = &map[string]*string{}
			nameCounters[fileCode.Path] = nameCounter
		}

		// Determine if this FileCode is the main one
		isMain := fileCode.Path == originalPath

		functionCounter := 0
		for _, fun := range fileCode.Functions {
			// Skip the main function
			if isMain && fun.Name == "main" {
				continue
			}

			counter := fmt.Sprintf("f__%d_%d", fileCodeCount, functionCounter)
			(*nameCounter)[fun.Name] = &counter
			functionCounter++
		}

		modCounter := 0
		for _, mod := range fileCode.Modules {
			modulePropCounters[mod] = util.NewOrderedMap[string, *string]()
			formattedName := fmt.Sprintf("m__%d_%d", fileCodeCount, modCounter)
			(*nameCounter)[mod.Name] = &formattedName
			modCounter++

			// Also precompute all properties and methods
			propCounter := 0
			propCounters := modulePropCounters[mod]

			for name := range mod.Declarations {
				counterFormatted := strconv.Itoa(propCounter)
				propCounters.Set(name, &counterFormatted)
				propCounter++
			}

			// Reset the counter
			propCounter = 0
			for name := range mod.Functions {
				methodName := fmt.Sprintf("%s__m_%d", formattedName, propCounter)
				propCounters.Set(name, &methodName)
				propCounter++
			}
		}

		fileCodeCount++
	}

	fileCodeCount = 0
	for _, fileCode := range fileCodes {
		// Skip the file if it has an external implementation
		if externalImpl[fileCode.Path] {
			continue
		}

		// Add the imports to the local counters
		for _, importPath := range fileCode.Imports {
			counters := nameCounters[importPath]

			if counters == nil {
				continue
			}

			for name, counter := range *counters {
				(*nameCounters[fileCode.Path])[name] = counter
			}
		}

		fileName := util.FileName(&fileCode.Path)

		// Emit a building state
		state.Emit(state.Building, fileName)

		// Determine if this FileCode is the main one
		isMain := fileCode.Path == originalPath

		fileIr := ir.BuildIr(
			fileCode,
			fileCodesMap,
			fileCodeCount,
			isMain,
			&originalPath,
			&traceCounters,
			&traceStrings,
			&usedStrings,
			&modulePropCounters,
			*nameCounters[fileCode.Path],
		)
		// Write the IR to the global builder
		globalBuilder.WriteString(fileIr)
		state.PassAllSpinners()
		fileCodeCount++
	}

	// Get the out directory path
	outDir := path.Join(dirPath, "out")

	// Make sure the output directory exists
	if !util.DirExists(outDir) {
		err := os.Mkdir(outDir, os.ModePerm)

		if err != nil {
			logger.Error("Could not create the output directory.")
			os.Exit(1)
		}
	}

	// Write the global IR to a file
	globalIrPath := path.Join(outDir, "program.flc")
	outPath := path.Join(outDir, "out")

	// Add an .exe extension if the user is on Windows
	if !isPOSIX {
		outPath += ".exe"
	}

	// Use a final builder to write the string references first
	finalBuilder := strings.Builder{}

	for address, str := range traceStrings {
		finalBuilder.WriteString("ref ")
		finalBuilder.WriteString(address)
		finalBuilder.WriteString(" str \"")
		finalBuilder.WriteString(*str)
		finalBuilder.WriteString("\"\n")
	}

	for str, address := range usedStrings.Storage {
		finalBuilder.WriteString("ref ")
		finalBuilder.WriteString(address)
		finalBuilder.WriteString(" str \"")
		finalBuilder.WriteString(str)
		finalBuilder.WriteString("\"\n")
	}

	for num, address := range traceCounters.Storage {
		finalBuilder.WriteString("ref ")
		finalBuilder.WriteString(address)
		finalBuilder.WriteString(" num ")
		finalBuilder.WriteString(strconv.Itoa(num))
		finalBuilder.WriteString("\n")
	}

	// Write the global builder
	finalBuilder.WriteString(globalBuilder.String())

	err = os.WriteFile(globalIrPath, []byte(finalBuilder.String()), os.ModePerm)

	if err != nil {
		logger.Error("Could not write the Fluent IR to a file.")
		os.Exit(1)
	}

	// If we are not supposed to call the compiler, return the IR path directly
	if !callCompiler {
		return globalIrPath
	}

	fmt.Println(ansi.Colorize(ansi.BoldBrightYellow, "⚠️ Invoking fluentc backend...."))
	fmt.Println(ansi.Colorize(ansi.BrightBlack, "⚠️ The output you will see from now on is coming from the fluentc command."))

	// Invoke the fluentc backend
	cmd = exec.Command("fluentc", "-O", strconv.Itoa(optimizationLevel), "-o", outPath, "-p", globalIrPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Errors are handled by the compiler backend
	_ = cmd.Run()

	return outPath
}
