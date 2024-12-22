package code

import (
	"surf/concurrent"
	"surf/object"
	"surf/token"
)

// SurfMod represents a Surf module
// a Surf module is somewhat similar to a class in OOP
type SurfMod struct {
	varDeclarations [][]token.Token
	properties      *concurrent.TypedConcurrentMap[string, object.SurfObject]
	publicMethods   map[string]*Function
	privateMethods  map[string]*Function
	name            string
	file            string
	public          bool
}

// NewSurfMod creates a new Surf module
func NewSurfMod(
	properties *concurrent.TypedConcurrentMap[string, object.SurfObject],
	publicMethods map[string]*Function,
	privateMethods map[string]*Function,
	name string,
	file string,
	varDeclarations [][]token.Token,
	public bool,
) SurfMod {
	return SurfMod{
		properties:      properties,
		publicMethods:   publicMethods,
		privateMethods:  privateMethods,
		name:            name,
		file:            file,
		varDeclarations: varDeclarations,
		public:          public,
	}
}

// FindMod finds a module in the given map
// and returns the module alongside a boolean
// indicating if the module was found
func FindMod(mods *map[string]map[string]*SurfMod, name string, file string) (*SurfMod, bool, bool) {
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

	return &SurfMod{}, false, false
}

// GetProperty returns the property with the given name
// alongside a boolean indicating if the property was found
func (sm *SurfMod) GetProperty(name string) (*object.SurfObject, bool) {
	prop, found := sm.properties.Load(name)
	return prop, found
}

// SetProperty sets the property with the given name
func (sm *SurfMod) SetProperty(name string, value object.SurfObject) {
	sm.properties.Store(name, value)
}

// GetMethod returns the method with the given name
// alongside a boolean indicating if the method was found
// and a boolean indicating if the method is public
func (sm *SurfMod) GetMethod(name string) (*Function, bool, bool) {
	method, found := sm.publicMethods[name]
	if found {
		return method, true, true
	}

	method, found = sm.privateMethods[name]
	return method, found, false
}

// GetName returns the name of the module
func (sm *SurfMod) GetName() string {
	return sm.name
}

// GetFile returns the file in which the module was defined
func (sm *SurfMod) GetFile() string {
	return sm.file
}

// GetVarDeclarations returns the variable declarations in the module
func (sm *SurfMod) GetVarDeclarations() [][]token.Token {
	return sm.varDeclarations
}

// IsPublic checks if the module is public
func (sm *SurfMod) IsPublic() bool {
	return sm.public
}
