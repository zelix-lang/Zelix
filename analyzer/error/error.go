/*
   The Fluent Programming Language
   -----------------------------------------------------
   Copyright (c) 2025 Rodrigo R. & All Fluent Contributors
   This program comes with ABSOLUTELY NO WARRANTY.
   For details type `fluent -l`. This is free software,
   and you are welcome to redistribute it under certain
   conditions; type `fluent -l -f` for details.
*/

package error

type Code int

const (
	Nothing                  Code = iota
	NameShouldBeSnakeCase         // E0001
	ParamTypeNothing              // E0002
	DataOutlivesStack             // E0003
	MustReturnAValue              // E0004
	UndefinedReference            // E0005
	UnusedVariable                // E0006
	Redefinition                  // E0007
	InvalidDereference            // E0008
	TypeMismatch                  // E0009
	ParameterCountMismatch        // E0010
	CannotInferType               // E0011
	ShouldNotReturn               // E0012
	CannotTakeAddress             // E0013
	InvalidPropAccess             // E0014
	IllegalPropAccess             // E0015
	ConstantReassignment          // E0016
	DoesNotHaveConstructor        // E0017
	ShouldNotHaveGenerics         // E0018
	ValueNotAssigned              // E0019
	CircularModuleDependency      // E0020
)

// Error represents an error with details about its location and additional information.
type Error struct {
	Line       int      // Line number where the error occurred.
	Column     int      // Column number where the error occurred.
	Code       Code     // The error code.
	Additional []string // Additional information about the error.
}
