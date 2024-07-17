package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/compiler"
	"golox/internal/debug"
	"golox/internal/opcode"
	"golox/internal/value"
)

type InterpretResult = uint8

const DEBUG_TRACE_EXECUTION bool = true

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
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

func (vm *VM) Interpret(source *[]byte) InterpretResult {
	compiler.Compile(source)
	return INTERPRET_OK
}

func (vm *VM) run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("          ")
			for slot := 0; slot < vm.stackTop; slot++ {
				fmt.Printf("[ ")
				vm.stack[slot].Print()
				fmt.Printf(" ]")
			}
			fmt.Printf("\n")
			debug.DisassembleInstruction(*vm.chunk, int(vm.chunk.Code[vm.ip]))
		}
		switch instruction := vm.readByte(); instruction {
		case opcode.Constant:
			constant := vm.readConstant()
			vm.push(constant)
		case opcode.ConstantLong:
			constant := vm.readConstant()
			vm.push(constant)
		case opcode.Add:
			vm.binaryOP(opcode.Add)
		case opcode.Subtract:
			vm.binaryOP(opcode.Subtract)
		case opcode.Multiply:
			vm.binaryOP(opcode.Multiply)
		case opcode.Divide:
			vm.binaryOP(opcode.Divide)
		case opcode.Negate:
			value := &vm.stack[vm.stackTop-1]
			*value = -*value
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

func (vm *VM) binaryOP(operator byte) {
	b := vm.pop()
	a := &vm.stack[vm.stackTop-1]
	switch operator {
	case opcode.Add:
		*a = *a + b
	case opcode.Subtract:
		*a = *a - b
	case opcode.Multiply:
		*a = *a * b
	case opcode.Divide:
		*a = *a / b
	default:
		err := fmt.Sprintf("Invalid binary operator %v", operator)
		panic(err)
	}
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}
