package types

type ZyroObjectType int

const (
	BooleanType ZyroObjectType = iota
	StringType
	IntType
	DecimalType
	NothingType
	ModType
)
