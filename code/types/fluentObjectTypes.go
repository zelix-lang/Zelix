package types

type FluentObjectType int

const (
	BooleanType FluentObjectType = iota
	StringType
	IntType
	DecimalType
	NothingType
	ModType
)
