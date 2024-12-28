package mod

import (
	"fluent/code"
	"fluent/code/types"
	"fluent/code/wrapper"
	"fluent/stack"
	"fluent/token"
	"fluent/tokenUtil/generic"
)

// FluentMod represents a Fluent module
// a Fluent module is somewhat similar to a class in OOP
type FluentMod struct {
	varDeclarations [][]token.Token
	properties      map[string]*wrapper.FluentObject
	methods         map[string]*code.Function
	name            string
	file            string
	public          bool
	trace           token.Token
	templates       []wrapper.TypeWrapper
	variables       *stack.Stack
	initialized     bool
}

// NewFluentMod creates a new Fluent module
func NewFluentMod(
	properties map[string]*wrapper.FluentObject,
	publicMethods map[string]*code.Function,
	privateMethods map[string]*code.Function,
	name string,
	file string,
	varDeclarations [][]token.Token,
	public bool,
	trace token.Token,
	templates []wrapper.TypeWrapper,
) FluentMod {
	mod := FluentMod{
		properties:      properties,
		methods:         make(map[string]*code.Function),
		name:            name,
		file:            file,
		varDeclarations: varDeclarations,
		public:          public,
		trace:           trace,
		templates:       templates,
		variables:       stack.NewStack(),
	}

	for key, value := range publicMethods {
		mod.methods[key] = value
	}

	for key, value := range privateMethods {
		mod.methods[key] = value
	}

	return mod
}

// FindMod finds a module in the given map
// and returns the module alongside a boolean
// indicating if the module was found and
// a boolean indicating if the module was found
// in the same file
func FindMod(mods *map[string]map[string]*FluentMod, name string, file string) (*FluentMod, bool, bool) {
	mod, found := (*mods)[file][name]
	if found {
		return mod, true, true
	}

	for _, fileMods := range *mods {
		mod, found = fileMods[name]

		if found {
			return mod, true, false
		}
	}

	return &FluentMod{}, false, false
}

// GetProperty returns the property with the given name
// alongside a boolean indicating if the property was found
func (sm *FluentMod) GetProperty(name string) (*wrapper.FluentObject, bool) {
	prop, found := sm.properties[name]
	return prop, found
}

// SetProperty sets the property with the given name
func (sm *FluentMod) SetProperty(name string, value wrapper.FluentObject) {
	sm.properties[name] = &value
}

// GetMethod returns the method with the given name
// alongside a boolean indicating if the method was found
// and a boolean indicating if the method is public
func (sm *FluentMod) GetMethod(name string) (*code.Function, bool, bool) {
	method, found := sm.methods[name]
	if found {
		return method, true, method.IsPublic()
	}

	return &code.Function{}, false, false
}

// GetName returns the name of the module
func (sm *FluentMod) GetName() string {
	return sm.name
}

// GetFile returns the file in which the module was defined
func (sm *FluentMod) GetFile() string {
	return sm.file
}

// GetVarDeclarations returns the variable declarations in the module
func (sm *FluentMod) GetVarDeclarations() [][]token.Token {
	return sm.varDeclarations
}

// IsPublic checks if the module is public
func (sm *FluentMod) IsPublic() bool {
	return sm.public
}

// GetTrace returns the trace of the module
func (sm *FluentMod) GetTrace() token.Token {
	return sm.trace
}

// GetMethods returns the methods of the module
func (sm *FluentMod) GetMethods() map[string]*code.Function {
	return sm.methods
}

// GetTemplates returns the templates of the module
func (sm *FluentMod) GetTemplates() []wrapper.TypeWrapper {
	return sm.templates
}

// BuildDummyWrapper builds a dummy type wrapper
// that represents the module
func (sm *FluentMod) BuildDummyWrapper() wrapper.TypeWrapper {
	return wrapper.ForceNewTypeWrapper(
		sm.name,
		sm.templates,
		types.ModType,
	)
}

// GetVariables returns the variables stack of the module
func (sm *FluentMod) GetVariables() *stack.Stack {
	return sm.variables
}

// BuildWithoutGenerics builds a new module, replacing
// generics with the given types
func (sm *FluentMod) BuildWithoutGenerics(types map[string]wrapper.TypeWrapper) FluentMod {
	varDeclarations := make([][]token.Token, 0)
	templates := make([]wrapper.TypeWrapper, len(sm.templates))
	properties := make(map[string]*wrapper.FluentObject)
	publicMethods := make(map[string]*code.Function)
	privateMethods := make(map[string]*code.Function)

	for i, template := range sm.templates {
		templates[i] = generic.ConvertGeneric(template, types)
	}

	for key, value := range sm.properties {
		newValue := wrapper.NewFluentObject(
			generic.ConvertGeneric(value.GetType(), types),
			value.GetValue(),
		)
		properties[key] = &newValue
	}

	for _, declaration := range sm.varDeclarations {
		varDeclarations = append(
			varDeclarations,
			generic.ConvertVariableGenerics(declaration, types),
		)
	}

	for key, value := range sm.methods {
		newFunction := value.BuildWithoutGenerics(types)

		if value.IsPublic() {
			publicMethods[key] = &newFunction
		} else {
			privateMethods[key] = &newFunction
		}
	}

	return NewFluentMod(
		properties,
		publicMethods,
		privateMethods,
		sm.name,
		sm.file,
		varDeclarations,
		sm.public,
		sm.trace,
		templates,
	)
}
