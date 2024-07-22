package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/compiler"
	"golox/internal/debug"
	"golox/internal/opcode"
	"golox/internal/value"
	"os"
)

type InterpretResult = uint8

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
	INTERPRET_NO_RESULT
)

const STACK_MAX int = 256

type VM struct {
	stack    [STACK_MAX]value.Value
	chunk    *chunk.Chunk
	ip       int
	stackTop int
	objects  *value.Obj
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Free() {
	// Implement cleanup logic
}

func (vm *VM) push(value value.Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() value.Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) peek(distance int) *value.Value {
	return &vm.stack[vm.stackTop-1-distance]
}

func (vm *VM) Interpret(source *[]byte) InterpretResult {
	c := chunk.NewChunk()

	if !compiler.Compile(source, c) {
		c.Free()
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = c
	vm.ip = 0

	result := vm.run()

	c.Free()
	return result
}

func (vm *VM) run() InterpretResult {
	for {
		if debug.TraceExecution {
			fmt.Printf("          ")
			for slot := 0; slot < vm.stackTop; slot++ {
				fmt.Printf("[ ")
				vm.stack[slot].Print()
				fmt.Printf(" ]")
			}
			fmt.Printf("\n")
			debug.DisassembleInstruction(vm.chunk, vm.ip)
		}

		switch instruction := vm.readByte(); instruction {
		case opcode.Constant, opcode.ConstantLong:
			constant := vm.readConstant()
			vm.push(constant)
		case opcode.Nil:
			vm.push(value.NewNil())
		case opcode.True:
			vm.push(value.NewBool(true))
		case opcode.False:
			vm.push(value.NewBool(false))
		case opcode.Equal:
			a := vm.peek(1)
			b := vm.pop()
			*a = value.NewBool((a.IsEqual(&b)))
		case opcode.NotEqual:
			a := vm.peek(1)
			b := vm.pop()
			*a = value.NewBool((!a.IsEqual(&b)))
		case opcode.Add:
			a := vm.peek(1)
			if b := vm.peek(0); a.IsString() && b.IsString() {
				b := vm.pop()
				*a = value.NewObjString(a.AsGoString() + b.AsGoString())
			} else if a.IsNumber() && b.IsNumber() {
				b := vm.pop()
				*a = value.NewNumber(a.AsNumber() + b.AsNumber())
			} else {
				vm.runtimeError(
					"Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
		case opcode.Greater, opcode.GreaterEqual, opcode.Less, opcode.LessEqual,
			opcode.Subtract, opcode.Multiply, opcode.Divide, opcode.Modulo:
			result := vm.binaryOP(instruction)
			if result != INTERPRET_NO_RESULT {
				return result
			}
		case opcode.Not:
			val := vm.peek(0)
			*val = value.NewBool(val.IsFalsey())
		case opcode.Negate:
			if val := vm.peek(0); !val.IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			} else {
				*val = value.NewNumber(-val.AsNumber())
			}
		case opcode.Return:
			fmt.Printf("====%v====\n", vm.chunk.Constants.Count())
			vm.pop().Print()
			fmt.Printf("\n")
			return INTERPRET_OK
		default:
			err := fmt.Sprintf("Unknown instruction %v", instruction)
			panic(err)
		}
	}
}

func (vm *VM) readByte() byte {
	defer func() { vm.ip++ }()
	return vm.chunk.Code[vm.ip]
}

func (vm *VM) readConstant() value.Value {
	return vm.chunk.Constants[vm.readByte()]
}

func (vm *VM) binaryOP(operator byte) InterpretResult {
	a := vm.peek(1)
	if b := vm.peek(0); !b.IsNumber() || !a.IsNumber() {
		vm.runtimeError("Operands must be numbers.")
		return INTERPRET_RUNTIME_ERROR
	}
	b := vm.pop()

	switch operator {
	case opcode.Greater:
		*a = value.NewBool(a.AsNumber() > b.AsNumber())
	case opcode.GreaterEqual:
		*a = value.NewBool(a.AsNumber() >= b.AsNumber())
	case opcode.Less:
		*a = value.NewBool(a.AsNumber() < b.AsNumber())
	case opcode.LessEqual:
		*a = value.NewBool(a.AsNumber() <= b.AsNumber())
	//case opcode.Add:
	//	*a = value.NewNumber(a.AsNumber() + b.AsNumber())
	case opcode.Subtract:
		*a = value.NewNumber(a.AsNumber() - b.AsNumber())
	case opcode.Multiply:
		*a = value.NewNumber(a.AsNumber() * b.AsNumber())
	case opcode.Divide:
		*a = value.NewNumber(a.AsNumber() / b.AsNumber())
	case opcode.Modulo:
		*a = value.NewNumber(float64(int(a.AsNumber()) % int(b.AsNumber())))
	default:
		err := fmt.Sprintf("Invalid binary operator %v", operator)
		panic(err)
	}

	return INTERPRET_NO_RESULT
}

func (vm *VM) runtimeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", vm.chunk.GetLine(vm.ip-1))
	vm.resetStack()
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}
