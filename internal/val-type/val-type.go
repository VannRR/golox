package valtype

const (
	Bool ValType = iota
	Nil
	Number
	Obj
)

type ValType = uint8
