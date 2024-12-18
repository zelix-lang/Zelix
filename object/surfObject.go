package object

type SurfObjectType int

const (
	BooleanType SurfObjectType = iota
	StringType
	IntType
	DecimalType
	NothingType
	ModType
)

type SurfObject struct {
	_type SurfObjectType
	value interface{}
}

func NewSurfObject(_type SurfObjectType, value any) SurfObject {
	return SurfObject{
		_type: _type,
		value: value,
	}
}

func (so *SurfObject) GetType() SurfObjectType {
	return so._type
}

func (so *SurfObject) GetValue() any {
	return so.value
}
