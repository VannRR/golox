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

func Name(op byte) string {
	switch op {
	case Constant:
		return "OpConstant"
	case ConstantLong:
		return "OpConstantLong"
	case Nil:
		return "OpNil"
	case True:
		return "OpTrue"
	case False:
		return "OpFalse"
	case Pop:
		return "OpPop"
	case GetLocal:
		return "OpGetLocal"
	case GetLocalLong:
		return "OpGetLocalLong"
	case SetLocal:
		return "OpSetLocal"
	case SetLocalLong:
		return "OpSetLocalLong"
	case GetGlobal:
		return "OpGetGlobal"
	case GetGlobalLong:
		return "OpGetGlobalLong"
	case DefineGlobal:
		return "OpDefineGlobal"
	case DefineGlobalLong:
		return "OpDefineGlobalLong"
	case SetGlobal:
		return "OpSetGlobal"
	case SetGlobalLong:
		return "OpSetGlobalLong"
	case Equal:
		return "OpEqual"
	case NotEqual:
		return "OpNotEqual"
	case Greater:
		return "OpGreater"
	case GreaterEqual:
		return "OpGreaterEqual"
	case Less:
		return "OpLess"
	case LessEqual:
		return "OpLessEqual"
	case Add:
		return "OpAdd"
	case Subtract:
		return "OpSubtract"
	case Multiply:
		return "OpMultiply"
	case Divide:
		return "OpDivide"
	case Not:
		return "OpNot"
	case Modulo:
		return "OpModulo"
	case Negate:
		return "OpNegate"
	case Print:
		return "OpPrint"
	case Jump:
		return "OpJump"
	case JumpIfFalse:
		return "OpJumpIfFalse"
	case Loop:
		return "OpLoop"
	case Return:
		return "OpReturn"
	default:
		return "OpUnknown"
	}
}
