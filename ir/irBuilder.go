package ir

import (
	"fluent/ast"
	"fluent/code"
	"fluent/ir/engine"
	"fluent/ir/wrapper"
	"os"
	"path"
	"strconv"
	"strings"
)

// The standard library's path
var stdPath = os.Getenv("FLUENT_STDLIB_PATH")

// BuildRuntimeInstructions locates all runtime functions and adds runtime instructions
// to the IR
func buildRuntimeInstructions(
	function *code.Function,
	file string,
	ir *wrapper.IrWrapper,
) {
	// Check if this file has any implementation in the "include" directory
	rawPath := strings.TrimPrefix(file, stdPath)

	// Delete slashes from the beginning of the path
	rawPath = strings.TrimPrefix(rawPath, "/")
	rawNoSuffix := strings.TrimSuffix(rawPath, ".fluent")
	includePath := path.Join(
		stdPath,
		"include",
		rawNoSuffix+".c",
	)

	// See if a C implementation exists
	if _, err := os.Stat(includePath); err != nil {
		return
	}

	// Add runtime instructions to the IR
	ir.AddRuntimeFunction(rawNoSuffix, function)
}

// EmitIR emits Fluent IR from the given file code
func EmitIR(fileCode ast.FileCode) string {
	// The resultant wrapper
	ir := wrapper.NewIrWrapper()

	// A global counter for variables and functions
	counter := 0

	// Iterate over all mods and compute the counter
	totalMods := fileCode.GetModules()

	for _, mods := range *totalMods {
		for _, mod := range mods {
			ir.AddMod(mod)

			// Iterate over all functions and compute the counter
			for _, function := range *mod.GetMethods() {
				// Increment the counter
				counter++

				// Save the function to the IR
				ir.AddFunction("x"+strconv.Itoa(counter), function)
			}
		}
	}

	// Iterate over all functions and compute the counter
	totalFunctions := fileCode.GetFunctions()
	for file, functions := range *totalFunctions {
		for _, function := range functions {
			// Check if the file comes from the stdlib
			if strings.HasPrefix(file, stdPath) {
				buildRuntimeInstructions(function, file, ir)
				continue
			}

			// Increment the counter
			counter++

			// Save the function to the IR
			ir.AddFunction("x"+strconv.Itoa(counter), function)
		}
	}

	return engine.MarshalIrWrapper(ir, &fileCode, counter)
}
