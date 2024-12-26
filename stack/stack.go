package stack

import (
	"fluent/code/wrapper"
)

// Stack represents a stack of variables for a function
// it is concurrent by design
type Stack struct {
	internal []map[string]*wrapper.FluentVariable
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		internal: make([]map[string]*wrapper.FluentVariable, 0),
	}
}

// CreateScope creates a new scope in the stack
func (s *Stack) CreateScope() {
	s.internal = append(s.internal, make(map[string]*wrapper.FluentVariable))
}

// Append appends a variable to the current scope
func (s *Stack) Append(key string, value wrapper.FluentObject, constant bool) {
	if len(s.internal) == 0 {
		s.CreateScope()
	}

	scope := s.internal[len(s.internal)-1]
	newVar := wrapper.NewFluentVariable(constant, value)
	scope[key] = &newVar

	// Update the scope in the stack
	s.internal[len(s.internal)-1] = scope
}

// DestroyScope destroys the current scope
func (s *Stack) DestroyScope() {
	if len(s.internal) == 0 {
		return
	}

	s.internal = s.internal[:len(s.internal)-1]
}

// Load retrieves a variable from the stack
func (s *Stack) Load(key string) (*wrapper.FluentVariable, bool) {
	for _, scope := range s.internal {
		value, found := scope[key]

		if found {
			return value, true
		}

	}

	dummyVar := wrapper.NewFluentVariable(false, wrapper.FluentObject{})
	return &dummyVar, false
}
