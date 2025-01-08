package stack

import (
	"fluent/code/wrapper"
	"fluent/logger"
	"fluent/token"
	"strings"
)

// Stack represents a stack of variables for a function
// it is concurrent by design
type Stack struct {
	internal   []map[string]*wrapper.FluentVariable
	loadedVars map[string]bool
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		internal: make([]map[string]*wrapper.FluentVariable, 0),
		// Used to keep track of the variables that have been
		// loaded (used), if the scope is destroyed and a
		// variable is not in this map, it means that the
		// variable was not used
		// Use a map for O(1) lookup
		loadedVars: make(map[string]bool),
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
func (s *Stack) DestroyScope(trace token.Token, inMod bool) {
	if len(s.internal) == 0 {
		return
	}

	lastScope := s.internal[len(s.internal)-1]

	for key := range lastScope {
		if strings.HasPrefix(key, "_") {
			// Variables that start with an underscore
			// are not checked for usage
			continue
		}

		_, found := s.loadedVars[key]
		if !found && !inMod && key != "this" {
			logger.TokenWarning(
				trace,
				"Unused variable "+key,
				"The variable "+key+" was declared but never used",
				"Remove the variable declaration",
				"Or add an underscore to the variable name to ignore this warning",
			)
		} else if found {
			delete(s.loadedVars, key)
		}
	}

	s.internal = s.internal[:len(s.internal)-1]
}

// Load retrieves a variable from the stack
func (s *Stack) Load(key string) (*wrapper.FluentVariable, bool) {
	for _, scope := range s.internal {
		value, found := scope[key]

		if found {
			// Mark the variable as used
			s.loadedVars[key] = true
			return value, true
		}

	}

	dummyVar := wrapper.NewFluentVariable(false, wrapper.FluentObject{})
	return &dummyVar, false
}
