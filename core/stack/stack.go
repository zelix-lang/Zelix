package stack

import (
	"zyro/code"
	"zyro/concurrent"
	"zyro/object"
)

// Stack represents a stack of variables for a function
// it is concurrent by design
type Stack struct {
	internal concurrent.TypedConcurrentSlice[concurrent.TypedConcurrentMap[string, code.ZyroVariable]]
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		internal: *concurrent.NewTypedConcurrentSlice[concurrent.TypedConcurrentMap[string, code.ZyroVariable]](),
	}
}

// CreateScope creates a new scope in the stack
func (s *Stack) CreateScope() {
	s.internal.Append(concurrent.NewTypedConcurrentMap[string, code.ZyroVariable]())
}

// Append appends a variable to the current scope
func (s *Stack) Append(key string, value object.ZyroObject, constant bool) {
	scope, _ := s.internal.Get(s.internal.Length() - 1)
	scope.Store(key, code.NewZyroVariable(constant, value))

	// Update the scope in the stack
	s.internal.Remove(s.internal.Length() - 1)
	s.internal.Append(&scope)
}

// DestroyScope destroys the current scope
func (s *Stack) DestroyScope() {
	s.internal.Remove(s.internal.Length() - 1)
}

// Load retrieves a variable from the stack
func (s *Stack) Load(key string) (*code.ZyroVariable, bool) {
	for i := 0; i < s.internal.Length(); i++ {
		scope, _ := s.internal.Get(i)

		value, found := scope.Load(key)

		if found {
			return value, true
		}

	}

	dummyVar := code.NewZyroVariable(false, object.ZyroObject{})
	return &dummyVar, false
}
