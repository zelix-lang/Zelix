package object

type ZyroObjectType int

const (
	BooleanType ZyroObjectType = iota
	StringType
	IntType
	DecimalType
	NothingType
	ModType
)

type ZyroObject struct {
	_type ZyroObjectType
	value interface{}
}

func NewZyroObject(_type ZyroObjectType, value any) ZyroObject {
	return ZyroObject{
		_type: _type,
		value: value,
	}
}

func (so *ZyroObject) GetType() ZyroObjectType {
	return so._type
}

func (so *ZyroObject) GetValue() any {
	return so.value
}
