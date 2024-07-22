package opcode

const (
	Constant byte = iota
	ConstantLong
	Nil
	True
	False
	Equal
	NotEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	Add
	Subtract
	Multiply
	Divide
	Not
	Modulo
	Negate
	Return
)
