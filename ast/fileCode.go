package ast

import (
	"surf/code"
	"surf/logger"
)

// FileCode is a representation of the file source code
type FileCode struct {
	// functions holds the functions of the current file
	// and of all the imported files.
	functions map[string]map[string]*Function
}

// GetFunctions returns the functions of the FileCode
func (fc *FileCode) GetFunctions() *map[string]map[string]*Function {
	return &fc.functions
}

// AddFunction adds a function to the FileCode
func (fc *FileCode) AddFunction(trace code.Token, file string, name string, function Function) {
	functions, ok := fc.functions[file]

	// Make sure the file exists in the map
	if !ok {
		fc.functions[file] = make(map[string]*Function)
		functions, _ = fc.functions[file]
	}

	for key := range functions {
		if key == name {
			logger.TokenError(
				trace,
				"Redefinition of function "+name,
				"The function "+name+" has already been defined in this file",
				"Change the name of the function",
			)
		}
	}

	functions[name] = &function
}

// GetFunction returns a function from the FileCode
func (fc *FileCode) GetFunction(file string, name string) (*Function, bool) {
	// Make sure the file exists in the map
	functions, ok := fc.functions[file]

	if !ok {
		return &Function{}, false
	}

	function, _ok := functions[name]
	return function, _ok
}

// LocateFunction locates a function in the FileCode and returns it
// along with a boolean indicating if the function was found
// and another boolean indicating if the function was found in the same file
func LocateFunction(
	functions map[string]map[string]*Function,
	file string,
	name string,
) (*Function, bool, bool) {
	fileFunctions, ok := functions[file]
	// Make sure the file exists in the map
	if !ok {
		return &Function{}, false, false
	}

	// Check if the function is in the current file
	if fun, ok := fileFunctions[name]; ok {
		return fun, true, true
	}

	// Check if the function is in an imported file
	for _, value := range functions {
		if fun, ok := value[name]; ok {
			return fun, true, false
		}
	}

	return &Function{}, false, false
}
