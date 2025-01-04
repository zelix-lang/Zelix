package wrapper

import (
	"fluent/code"
	"fluent/code/mod"
	"fluent/code/wrapper"
)

// IrWrapper represents a wrapper for Fluent IR
type IrWrapper struct {
	// runtime holds a map of runtime functions categorized by their names
	runtime map[string][]string
	// functions holds a map of functions categorized by their names
	functions map[*code.Function]string
	// globalVars holds a map of global variables categorized by their names
	globalVars map[string]*wrapper.FluentObject
	// mods holds a map of FluentMod objects and the times each one of them has been built
	mods map[*mod.FluentMod]*int
}

// NewIrWrapper creates a new IrWrapper
func NewIrWrapper() *IrWrapper {
	return &IrWrapper{
		runtime:    make(map[string][]string),
		functions:  make(map[*code.Function]string),
		globalVars: make(map[string]*wrapper.FluentObject),
		mods:       make(map[*mod.FluentMod]*int),
	}
}

// AddFunction adds a function to the IrWrapper
func (ir *IrWrapper) AddFunction(name string, function *code.Function) {
	ir.functions[function] = name
}

// GetFunction gets a function from the IrWrapper
func (ir *IrWrapper) GetFunction(function *code.Function) string {
	return ir.functions[function]
}

// GetFunctions gets all functions from the IrWrapper
func (ir *IrWrapper) GetFunctions() map[*code.Function]string {
	return ir.functions
}

// AddRuntimeFunction adds a runtime function to the IrWrapper
func (ir *IrWrapper) AddRuntimeFunction(name string, function *code.Function) {
	if _, found := ir.runtime[name]; !found {
		ir.runtime[name] = make([]string, 0)
	}

	ir.runtime[name] = append(ir.runtime[name], function.GetName())
}

// GetRuntimeFunctions gets all runtime functions from the IrWrapper
func (ir *IrWrapper) GetRuntimeFunctions() map[string][]string {
	return ir.runtime
}

// AddGlobalVar adds a global variable to the IrWrapper
func (ir *IrWrapper) AddGlobalVar(name string, value *wrapper.FluentObject) {
	ir.globalVars[name] = value
}

// GetGlobalVar gets a global variable from the IrWrapper
func (ir *IrWrapper) GetGlobalVar(name string) *wrapper.FluentObject {
	return ir.globalVars[name]
}

// GetGlobalVars gets all global variables from the IrWrapper
func (ir *IrWrapper) GetGlobalVars() map[string]*wrapper.FluentObject {
	return ir.globalVars
}

// AddMod adds a mod to the IrWrapper
func (ir *IrWrapper) AddMod(mod *mod.FluentMod) {
	if _, found := ir.mods[mod]; !found {
		ir.mods[mod] = new(int)
	}

	*ir.mods[mod]++
}

// GetMods gets all mods from the IrWrapper
func (ir *IrWrapper) GetMods() map[*mod.FluentMod]*int {
	return ir.mods
}

// GetMod gets a mod from the IrWrapper
func (ir *IrWrapper) GetMod(mod *mod.FluentMod) *int {
	return ir.mods[mod]
}
