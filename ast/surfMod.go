package ast

import (
	"surf/code"
	"surf/concurrent"
	"surf/object"
)

// SurfMod represents a Surf module
// a Surf module is somewhat similar to a class in OOP
type SurfMod struct {
	varDeclarations [][]code.Token
	properties      *concurrent.TypedConcurrentMap[string, object.SurfObject]
	publicMethods   map[string]*Function
	privateMethods  map[string]*Function
	name            string
	file            string
}

// NewSurfMod creates a new Surf module
func NewSurfMod(
	properties *concurrent.TypedConcurrentMap[string, object.SurfObject],
	publicMethods map[string]*Function,
	privateMethods map[string]*Function,
	name string,
	file string,
	varDeclarations [][]code.Token,
) SurfMod {
	return SurfMod{
		properties:      properties,
		publicMethods:   publicMethods,
		privateMethods:  privateMethods,
		name:            name,
		file:            file,
		varDeclarations: varDeclarations,
	}
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
	return method, found, true
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
func (sm *SurfMod) GetVarDeclarations() [][]code.Token {
	return sm.varDeclarations
}
