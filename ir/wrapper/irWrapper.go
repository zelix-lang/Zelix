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
	// runtimeFunctions holds a map of runtime functions that allows quick lookup
	runtimeFunctions map[*code.Function]string
	// functions holds a map of functions categorized by their names
	functions map[*code.Function]string
	// globalVars holds a map of global variables categorized by their names
	globalVars map[string]*wrapper.FluentObject
	// mods holds a map of FluentMod objects and their counts
	mods map[*mod.FluentMod]int
	// modsProps holds a map of FluentMod properties and their counts
	modsProps map[*mod.FluentMod]map[string]int
}

// NewIrWrapper creates a new IrWrapper
func NewIrWrapper() *IrWrapper {
	return &IrWrapper{
		runtime:          make(map[string][]string),
		functions:        make(map[*code.Function]string),
		globalVars:       make(map[string]*wrapper.FluentObject),
		mods:             make(map[*mod.FluentMod]int),
		runtimeFunctions: make(map[*code.Function]string),
		modsProps:        make(map[*mod.FluentMod]map[string]int),
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

	ir.runtimeFunctions[function] = name
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
func (ir *IrWrapper) AddMod(mod *mod.FluentMod, counter int) {
	ir.mods[mod] = counter
}

// GetMods gets all mods from the IrWrapper
func (ir *IrWrapper) GetMods() map[*mod.FluentMod]int {
	return ir.mods
}

// GetMod gets a mod from the IrWrapper
func (ir *IrWrapper) GetMod(mod *mod.FluentMod) int {
	return ir.mods[mod]
}

// GetRuntimeFunction checks if a function is runtime and returns its path
func (ir *IrWrapper) GetRuntimeFunction(fun *code.Function) (string, bool) {
	function, found := ir.runtimeFunctions[fun]
	return function, found
}

// AddModProp adds a mod property to the IrWrapper
func (ir *IrWrapper) AddModProp(mod *mod.FluentMod, prop string, counter int) {
	if _, found := ir.modsProps[mod]; !found {
		ir.modsProps[mod] = make(map[string]int)
	}

	ir.modsProps[mod][prop] = counter
}

// GetModProp gets a mod property from the IrWrapper
func (ir *IrWrapper) GetModProp(mod *mod.FluentMod, prop string) int {
	return ir.modsProps[mod][prop]
}
