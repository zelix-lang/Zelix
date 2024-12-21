package stack

import (
	"surf/code"
	"surf/concurrent"
	"surf/object"
)

// Stack represents a stack of variables for a function
// it is concurrent by design
type Stack struct {
	internal concurrent.TypedConcurrentSlice[concurrent.TypedConcurrentMap[string, code.SurfVariable]]
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		internal: *concurrent.NewTypedConcurrentSlice[concurrent.TypedConcurrentMap[string, code.SurfVariable]](),
	}
}

// CreateScope creates a new scope in the stack
func (s *Stack) CreateScope() {
	s.internal.Append(concurrent.NewTypedConcurrentMap[string, code.SurfVariable]())
}

// Append appends a variable to the current scope
func (s *Stack) Append(key string, value object.SurfObject) {
	scope, _ := s.internal.Get(s.internal.Length() - 1)
	scope.Store(key, code.NewSurfVariable(false, value))

	// Update the scope in the stack
	s.internal.Remove(s.internal.Length() - 1)
	s.internal.Append(&scope)
}

// DestroyScope destroys the current scope
func (s *Stack) DestroyScope() {
	s.internal.Remove(s.internal.Length() - 1)
}

// Load retrieves a variable from the stack
func (s *Stack) Load(key string) (object.SurfObject, bool) {
	for i := 0; i < s.internal.Length(); i++ {
		scope, _ := s.internal.Get(i)

		value, found := scope.Load(key)

		if found {
			return value.GetValue(), true
		}

	}

	return object.SurfObject{}, false
}
