package wrapper

// FluentObject is a struct that wraps a type and a value.
type FluentObject struct {
	_type TypeWrapper // _type represents the type of the wrapped value.
	value interface{} // value holds the actual value being wrapped.
}

// NewFluentObject creates a new FluentObject with the given type and value.
// Parameters:
//   - _type: The type of the value being wrapped.
//   - value: The actual value to be wrapped.
//
// Returns:
//
//	A new instance of FluentObject.
func NewFluentObject(_type TypeWrapper, value any) FluentObject {
	return FluentObject{
		_type: _type,
		value: value,
	}
}

// GetType returns the type of the wrapped value.
// Returns:
//
//	The type of the wrapped value.
func (so *FluentObject) GetType() TypeWrapper {
	return so._type
}

// GetValue returns the wrapped value.
// Returns:
//
//	The wrapped value.
func (so *FluentObject) GetValue() any {
	return so.value
}
