package opcode

const (
	Constant byte = iota
	ConstantLong
	Nil
	True
	False
	Pop
	GetLocal
	GetLocalLong
	SetLocal
	SetLocalLong
	GetGlobal
	GetGlobalLong
	DefineGlobal
	DefineGlobalLong
	SetGlobal
	SetGlobalLong
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
	Print
	Jump
	JumpIfFalse
	Loop
	Return
)

var Name = map[byte]string{
	Constant:         "OpConstant",
	ConstantLong:     "OpConstantLong",
	Nil:              "OpNil",
	True:             "OpTrue",
	False:            "OpFalse",
	Pop:              "OpPop",
	GetLocal:         "OpGetLocal",
	GetLocalLong:     "OpGetLocalLong",
	SetLocal:         "OpSetLocal",
	SetLocalLong:     "OpSetLocalLong",
	GetGlobal:        "OpGetGlobal",
	GetGlobalLong:    "OpGetGlobalLong",
	DefineGlobal:     "OpDefineGlobal",
	DefineGlobalLong: "OpDefineGlobalLong",
	SetGlobal:        "OpSetGlobal",
	SetGlobalLong:    "OpSetGlobalLong",
	Equal:            "OpEqual",
	NotEqual:         "OpNotEqual",
	Greater:          "OpGreater",
	GreaterEqual:     "OpGreaterEqual",
	Less:             "OpLess",
	LessEqual:        "OpLessEqual",
	Add:              "OpAdd",
	Subtract:         "OpSubtract",
	Multiply:         "OpMultiply",
	Divide:           "OpDivide",
	Not:              "OpNot",
	Modulo:           "OpModulo",
	Negate:           "OpNegate",
	Print:            "OpPrint",
	Jump:             "OpJump",
	JumpIfFalse:      "OpJumpIfFalse",
	Loop:             "OpLoop",
	Return:           "OpReturn",
}
