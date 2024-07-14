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

func (vm *VM) Interpret(c *Chunk) InterpretResult {
	vm.chunk = c
	vm.ip = 0
	return vm.run()
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
		case OP_RETURN:
			vm.pop().Print()
			fmt.Printf("\n")
			return INTERPRET_OK
		}
	}
}

func (vm *VM) readByte() uint8 {
	defer func() { vm.ip++ }()
	return vm.chunk.code[vm.ip]
}

func (vm *VM) readConstant() Value {
	return vm.chunk.constants[vm.readByte()]
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}
