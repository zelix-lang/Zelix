package stack

import (
	"surf/concurrent"
	"surf/object"
)

// StaticStack represents a stack of variables for a function
// it is concurrent by design
type StaticStack struct {
	internal concurrent.TypedConcurrentSlice[concurrent.TypedConcurrentMap[string, object.SurfObjectType]]
}

// NewStaticStack creates a new stack
func NewStaticStack() *StaticStack {
	return &StaticStack{
		internal: *concurrent.NewTypedConcurrentSlice[concurrent.TypedConcurrentMap[string, object.SurfObjectType]](),
	}
}

// CreateScope creates a new scope in the stack
func (s *StaticStack) CreateScope() {
	s.internal.Append(concurrent.NewTypedConcurrentMap[string, object.SurfObjectType]())
}

// Append appends a variable to the current scope
func (s *StaticStack) Append(key string, value object.SurfObjectType) {
	scope, _ := s.internal.Get(s.internal.Length() - 1)
	scope.Store(key, value)

	// Update the scope in the stack
	s.internal.Remove(s.internal.Length() - 1)
	s.internal.Append(&scope)
}

// DestroyScope destroys the current scope
func (s *StaticStack) DestroyScope() {
	s.internal.Remove(s.internal.Length() - 1)
}

// Load retrieves a variable from the stack
func (s *StaticStack) Load(key string) (object.SurfObjectType, bool) {
	for i := 0; i < s.internal.Length(); i++ {
		scope, _ := s.internal.Get(i)

		value, found := scope.Load(key)

		if found {
			return *value, true
		}

	}

	return object.NothingType, false
}
