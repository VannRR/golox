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
	chunk    *chunk.Chunk
	ip       int
	stack    [STACK_MAX]value.Value
	stackTop int
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
			vm.push(value.NewNilVal())
		case opcode.True:
			vm.push(value.NewBoolVal(true))
		case opcode.False:
			vm.push(value.NewBoolVal(false))
		case opcode.Add, opcode.Subtract, opcode.Multiply, opcode.Divide, opcode.Modulo:
			result := vm.binaryOP(instruction)
			if result != INTERPRET_NO_RESULT {
				return result
			}
		case opcode.Not:
			val := vm.peek(0)
			*val = value.NewBoolVal(val.IsFalsey())
		case opcode.Negate:
			if val := vm.peek(0); !val.IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			} else {
				*val = value.NewNumberVal(-val.AsNumber())
			}
		case opcode.Return:
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
	case opcode.Add:
		*a = value.NewNumberVal(a.AsNumber() + b.AsNumber())
	case opcode.Subtract:
		*a = value.NewNumberVal(a.AsNumber() - b.AsNumber())
	case opcode.Multiply:
		*a = value.NewNumberVal(a.AsNumber() * b.AsNumber())
	case opcode.Divide:
		*a = value.NewNumberVal(a.AsNumber() / b.AsNumber())
	case opcode.Modulo:
		*a = value.NewNumberVal(float64(int(a.AsNumber()) % int(b.AsNumber())))
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
