package wrapper

type ZyroObject struct {
	_type TypeWrapper
	value interface{}
}

func NewZyroObject(_type TypeWrapper, value any) ZyroObject {
	return ZyroObject{
		_type: _type,
		value: value,
	}
}

func (so *ZyroObject) GetType() TypeWrapper {
	return so._type
}

func (so *ZyroObject) GetValue() any {
	return so.value
}
