package main

import "fmt"

type InterpretResult = uint8

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const STACK_MAX int = 256

type VM struct {
	chunk    *Chunk
	ip       int
	stack    [STACK_MAX]Value
	stackTop int
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Free() {
	// Implement cleanup logic
}

func (vm *VM) push(value Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) Interpret(source *[]byte) InterpretResult {
	Compile(source)
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
			disassembleInstruction(*vm.chunk, int(vm.chunk.code[vm.ip]))
		}
		switch instruction := vm.readByte(); instruction {
		case OP_CONSTANT:
			constant := vm.readConstant()
			vm.push(constant)
		case OP_CONSTANT_LONG:
			constant := vm.readConstant()
			vm.push(constant)
		case OP_ADD:
			vm.binaryOP(OP_ADD)
		case OP_SUBTRACT:
			vm.binaryOP(OP_SUBTRACT)
		case OP_MULTIPLY:
			vm.binaryOP(OP_MULTIPLY)
		case OP_DIVIDE:
			vm.binaryOP(OP_DIVIDE)
		case OP_NEGATE:
			value := &vm.stack[vm.stackTop-1]
			*value = -*value
		case OP_RETURN:
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
	return vm.chunk.code[vm.ip]
}

func (vm *VM) readConstant() Value {
	return vm.chunk.constants[vm.readByte()]
}

func (vm *VM) binaryOP(operator byte) {
	b := vm.pop()
	a := &vm.stack[vm.stackTop-1]
	switch operator {
	case OP_ADD:
		*a = *a + b
	case OP_SUBTRACT:
		*a = *a - b
	case OP_MULTIPLY:
		*a = *a * b
	case OP_DIVIDE:
		*a = *a / b
	default:
		err := fmt.Sprintf("Invalid binary operator %v", operator)
		panic(err)
	}
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}
