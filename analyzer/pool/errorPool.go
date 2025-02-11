/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package pool

import error3 "fluent/analyzer/error"

// ErrorPool is a struct that holds a slice of errors and a count of errors.
type ErrorPool struct {
	Errors []error3.Error
	Count  int
}

// NewErrorPool creates and returns a new ErrorPool instance.
// Returns: *ErrorPool - a pointer to the newly created ErrorPool.
func NewErrorPool() *ErrorPool {
	return &ErrorPool{
		Errors: make([]error3.Error, 0),
		Count:  0,
	}
}

// AddError adds a new error to the ErrorPool if the count is less than 10.
// Parameters:
//   - err: error3.Error - the error to be added.
func (e *ErrorPool) AddError(err error3.Error) {
	// Keep a max of 10 errors
	if e.Count >= 10 {
		return
	}

	// Skip if the error is nothing
	if err.Code == error3.Nothing {
		return
	}

	e.Errors = append(e.Errors, err)
	e.Count++
}

// Extend adds multiple errors to the ErrorPool.
// Parameters:
//   - errs: []error3.Error - a slice of errors to be added.
func (e *ErrorPool) Extend(errs []error3.Error) {
	for _, err := range errs {
		e.AddError(err)
	}
}
