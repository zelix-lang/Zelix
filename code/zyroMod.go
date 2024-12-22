package code

import (
	"zyro/concurrent"
	"zyro/object"
	"zyro/token"
)

// ZyroMod represents a Zyro module
// a Zyro module is somewhat similar to a class in OOP
type ZyroMod struct {
	varDeclarations [][]token.Token
	properties      *concurrent.TypedConcurrentMap[string, object.ZyroObject]
	methods         map[string]*Function
	name            string
	file            string
	public          bool
	trace           token.Token
}

// NewZyroMod creates a new Zyro module
func NewZyroMod(
	properties *concurrent.TypedConcurrentMap[string, object.ZyroObject],
	publicMethods map[string]*Function,
	privateMethods map[string]*Function,
	name string,
	file string,
	varDeclarations [][]token.Token,
	public bool,
	trace token.Token,
) ZyroMod {
	mod := ZyroMod{
		properties:      properties,
		methods:         make(map[string]*Function),
		name:            name,
		file:            file,
		varDeclarations: varDeclarations,
		public:          public,
		trace:           trace,
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
// indicating if the module was found
func FindMod(mods *map[string]map[string]*ZyroMod, name string, file string) (*ZyroMod, bool, bool) {
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

	return &ZyroMod{}, false, false
}

// GetProperty returns the property with the given name
// alongside a boolean indicating if the property was found
func (sm *ZyroMod) GetProperty(name string) (*object.ZyroObject, bool) {
	prop, found := sm.properties.Load(name)
	return prop, found
}

// SetProperty sets the property with the given name
func (sm *ZyroMod) SetProperty(name string, value object.ZyroObject) {
	sm.properties.Store(name, value)
}

// GetMethod returns the method with the given name
// alongside a boolean indicating if the method was found
// and a boolean indicating if the method is public
func (sm *ZyroMod) GetMethod(name string) (*Function, bool, bool) {
	method, found := sm.methods[name]
	if found {
		return method, true, method.public
	}

	return &Function{}, false, false
}

// GetName returns the name of the module
func (sm *ZyroMod) GetName() string {
	return sm.name
}

// GetFile returns the file in which the module was defined
func (sm *ZyroMod) GetFile() string {
	return sm.file
}

// GetVarDeclarations returns the variable declarations in the module
func (sm *ZyroMod) GetVarDeclarations() [][]token.Token {
	return sm.varDeclarations
}

// IsPublic checks if the module is public
func (sm *ZyroMod) IsPublic() bool {
	return sm.public
}

// GetTrace returns the trace of the module
func (sm *ZyroMod) GetTrace() token.Token {
	return sm.trace
}

// GetMethods returns the methods of the module
func (sm *ZyroMod) GetMethods() map[string]*Function {
	return sm.methods
}
