package valtype

const (
	Bool Type = iota
	Nil
	Number
	ObjString
)

type Type = uint8
