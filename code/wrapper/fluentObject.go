package wrapper

type FluentObject struct {
	_type TypeWrapper
	value interface{}
}

func NewFluentObject(_type TypeWrapper, value any) FluentObject {
	return FluentObject{
		_type: _type,
		value: value,
	}
}

func (so *FluentObject) GetType() TypeWrapper {
	return so._type
}

func (so *FluentObject) GetValue() any {
	return so.value
}
