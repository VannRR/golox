package opcode

const (
	Constant byte = iota
	ConstantLong
	Nil
	True
	False
	Add
	Subtract
	Multiply
	Divide
	Not
	Modulo
	Negate
	Return
)
