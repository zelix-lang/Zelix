package stack

import (
	"zyro/code/wrapper"
	"zyro/concurrent"
)

// Stack represents a stack of variables for a function
// it is concurrent by design
type Stack struct {
	internal concurrent.TypedConcurrentSlice[concurrent.TypedConcurrentMap[string, wrapper.ZyroVariable]]
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		internal: *concurrent.NewTypedConcurrentSlice[concurrent.TypedConcurrentMap[string, wrapper.ZyroVariable]](),
	}
}

// CreateScope creates a new scope in the stack
func (s *Stack) CreateScope() {
	s.internal.Append(concurrent.NewTypedConcurrentMap[string, wrapper.ZyroVariable]())
}

// Append appends a variable to the current scope
func (s *Stack) Append(key string, value wrapper.ZyroObject, constant bool) {
	scope, _ := s.internal.Get(s.internal.Length() - 1)
	scope.Store(key, wrapper.NewZyroVariable(constant, value))

	// Update the scope in the stack
	s.internal.Remove(s.internal.Length() - 1)
	s.internal.Append(&scope)
}

// DestroyScope destroys the current scope
func (s *Stack) DestroyScope() {
	s.internal.Remove(s.internal.Length() - 1)
}

// Load retrieves a variable from the stack
func (s *Stack) Load(key string) (*wrapper.ZyroVariable, bool) {
	for i := 0; i < s.internal.Length(); i++ {
		scope, _ := s.internal.Get(i)

		value, found := scope.Load(key)

		if found {
			return value, true
		}

	}

	dummyVar := wrapper.NewZyroVariable(false, wrapper.ZyroObject{})
	return &dummyVar, false
}
