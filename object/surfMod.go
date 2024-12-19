package object

import (
	"surf/ast"
	"surf/concurrent"
)

// SurfMod represents a Surf module
// a Surf module is somewhat similar to a class in OOP
type SurfMod struct {
	// Properties of the module
	properties     *concurrent.TypedConcurrentMap[string, SurfObject]
	publicMethods  map[string]*ast.Function
	privateMethods map[string]*ast.Function
}

// NewSurfMod creates a new Surf module
func NewSurfMod(
	properties *concurrent.TypedConcurrentMap[string, SurfObject],
	publicMethods map[string]*ast.Function,
	privateMethods map[string]*ast.Function,
) *SurfMod {
	return &SurfMod{
		properties:     properties,
		publicMethods:  publicMethods,
		privateMethods: privateMethods,
	}
}

// GetProperty returns the property with the given name
// alongside a boolean indicating if the property was found
func (sm *SurfMod) GetProperty(name string) (*SurfObject, bool) {
	prop, found := sm.properties.Load(name)
	return prop, found
}

// SetProperty sets the property with the given name
func (sm *SurfMod) SetProperty(name string, value SurfObject) {
	sm.properties.Store(name, value)
}

// GetMethod returns the method with the given name
// alongside a boolean indicating if the method was found
// and a boolean indicating if the method is public
func (sm *SurfMod) GetMethod(name string) (*ast.Function, bool, bool) {
	method, found := sm.publicMethods[name]
	return method, found, true
}
