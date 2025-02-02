/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package stack

import "fluent/analyzer/variable"

// Stack represents a collection of variables.
type Stack struct {
	// Variables is a map where the key is a string and the value is a variable.Variable.
	Variables map[string]variable.Variable
	// UsedVariables is a map of strings representing the names of the variables that have been used.
	UsedVariables map[string]struct{}
}

// ScopedStack represents a stack with multiple scopes.
type ScopedStack struct {
	// Scopes is a slice of Stack, each representing a different scope.
	Scopes []Stack
}

// NewScope adds a new scope to the ScopedStack.
func (s *ScopedStack) NewScope() {
	s.Scopes = append(s.Scopes, Stack{
		Variables:     make(map[string]variable.Variable),
		UsedVariables: make(map[string]struct{}),
	})
}

// DestroyScope removes the last scope from the ScopedStack and returns a slice of unused variable names.
//
// Returns:
//
//	[]string: A slice of strings representing the names of unused variables in the last scope.
func (s *ScopedStack) DestroyScope() []*string {
	// Get the last scope
	lastScope := s.Scopes[len(s.Scopes)-1]

	s.Scopes = s.Scopes[:len(s.Scopes)-1]
	unusedVars := make([]*string, 0)

	// Iterate over the variables in the last scope
	for key := range lastScope.Variables {
		// Skip variables that are intended to be unused
		if key[0] == '_' {
			continue
		}

		// See if the variable was used
		if _, ok := lastScope.UsedVariables[key]; !ok {
			unusedVars = append(unusedVars, &key)
		}
	}

	// Return the unused variables
	return unusedVars
}

// Load retrieves a variable by name from the most recent scope where it exists.
// If the variable is found, it is marked as used.
//
// Parameters:
//
//	name (string): The name of the variable to retrieve.
//
// Returns:
//
//	*variable.Variable: A pointer to the variable if found, or nil if not found.
func (s *ScopedStack) Load(name *string) *variable.Variable {
	// Iterate over the scopes in reverse order
	for i := len(s.Scopes) - 1; i >= 0; i-- {
		scope := s.Scopes[i]

		if v, ok := scope.Variables[*name]; ok {
			scope.UsedVariables[*name] = struct{}{}
			return &v
		}
	}

	return nil
}

// Append adds a variable to the last scope in the ScopedStack.
//
// Parameters:
//
//	name (string): The name of the variable to add.
//	variable2 (variable.Variable): The variable to add.
func (s *ScopedStack) Append(name string, variable2 variable.Variable) {
	// Get the last scope
	lastScope := &s.Scopes[len(s.Scopes)-1]
	lastScope.Variables[name] = variable2
}
