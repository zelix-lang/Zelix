package wrapper

// FluentVariable represents a variable in a Fluent program
type FluentVariable struct {
	constant bool
	value    FluentObject
}

// NewFluentVariable creates a new Fluent variable
func NewFluentVariable(constant bool, value FluentObject) FluentVariable {
	return FluentVariable{
		constant: constant,
		value:    value,
	}
}

// IsConstant returns true if the variable is a constant
func (sv *FluentVariable) IsConstant() bool {
	return sv.constant
}

// GetValue returns the value of the variable
func (sv *FluentVariable) GetValue() FluentObject {
	return sv.value
}
