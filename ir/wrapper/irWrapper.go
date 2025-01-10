package wrapper

import (
	"fluent/code"
	"fluent/code/mod"
)

// IrWrapper represents a wrapper for Fluent IR
type IrWrapper struct {
	// runtime holds a map of runtime functions categorized by their names
	runtime map[string][]string
	// runtimeFunctions holds a map of runtime functions that allows quick lookup
	runtimeFunctions map[*code.Function]string
	// functions holds a map of functions categorized by their names
	functions map[*code.Function]string
	// mods holds a map of FluentMod objects and their counts
	mods map[*mod.FluentMod]int
	// modsProps holds a map of FluentMod properties and their counts
	modsProps map[string]map[string]int
	// genericMods holds a map of already-built generic mods
	genericMods map[string]*mod.FluentMod
}

// NewIrWrapper creates a new IrWrapper
func NewIrWrapper() *IrWrapper {
	return &IrWrapper{
		runtime:          make(map[string][]string),
		functions:        make(map[*code.Function]string),
		mods:             make(map[*mod.FluentMod]int),
		runtimeFunctions: make(map[*code.Function]string),
		modsProps:        make(map[string]map[string]int),
		genericMods:      make(map[string]*mod.FluentMod),
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

// AddMod adds a mod to the IrWrapper
func (ir *IrWrapper) AddMod(mod *mod.FluentMod, counter int) {
	ir.mods[mod] = counter
}

// GetMods gets all mods from the IrWrapper
func (ir *IrWrapper) GetMods() map[*mod.FluentMod]int {
	return ir.mods
}

// GetMod gets a mod from the IrWrapper
func (ir *IrWrapper) GetMod(mod *mod.FluentMod) (int, bool) {
	num, found := ir.mods[mod]

	if !found {
		return -1, false
	}

	return num, found
}

// GetRuntimeFunction checks if a function is runtime and returns its path
func (ir *IrWrapper) GetRuntimeFunction(fun *code.Function) (string, bool) {
	function, found := ir.runtimeFunctions[fun]
	return function, found
}

// AddModProp adds a mod property to the IrWrapper
func (ir *IrWrapper) AddModProp(mod *mod.FluentMod, prop string, counter int) {
	dummyWrapper := mod.BuildDummyWrapper()
	wrapperMarshal := dummyWrapper.Marshal()

	if _, found := ir.modsProps[wrapperMarshal]; !found {
		ir.modsProps[wrapperMarshal] = make(map[string]int)
	}

	ir.modsProps[wrapperMarshal][prop] = counter
}

// GetModProp gets a mod property from the IrWrapper
func (ir *IrWrapper) GetModProp(mod *mod.FluentMod, prop string) int {
	dummyWrapper := mod.BuildDummyWrapper()
	wrapperMarshal := dummyWrapper.Marshal()

	return ir.modsProps[wrapperMarshal][prop]
}

// AddGenericMod adds a generic mod to the IrWrapper
func (ir *IrWrapper) AddGenericMod(name string, mod *mod.FluentMod) {
	ir.genericMods[name] = mod
}

// GetGenericMod gets a generic mod from the IrWrapper
func (ir *IrWrapper) GetGenericMod(name string) (*mod.FluentMod, bool) {
	module, found := ir.genericMods[name]
	return module, found
}

// GetGenericMods gets all generic mods from the IrWrapper
func (ir *IrWrapper) GetGenericMods() map[string]*mod.FluentMod {
	return ir.genericMods
}
