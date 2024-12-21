package code

import "surf/object"

// SurfVariable represents a variable in a Surf program
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

// IsConstant returns true if the variable is a constant
func (sv *SurfVariable) IsConstant() bool {
	return sv.constant
}

// GetValue returns the value of the variable
func (sv *SurfVariable) GetValue() object.SurfObject {
	return sv.value
}
