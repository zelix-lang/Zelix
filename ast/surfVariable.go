package ast

import "surf/object"

// SurfVariable represents a variable in Surf
type SurfVariable struct {
	constant bool
	value    object.SurfObject
}

// NewSurfVariable creates a new Surf variable
func NewSurfVariable(constant bool, value object.SurfObject) SurfVariable {
	return SurfVariable{
		constant: constant,
		value:    value,
	}
}

// IsConstant checks if the variable is constant
func (sv *SurfVariable) IsConstant() bool {
	return sv.constant
}

// GetValue returns the value of the variable
func (sv *SurfVariable) GetValue() object.SurfObject {
	return sv.value
}

// SetValue sets the value of the variable
func (sv *SurfVariable) SetValue(value object.SurfObject) {
	sv.value = value
}
