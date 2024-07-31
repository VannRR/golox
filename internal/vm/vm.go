package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/common"
	"golox/internal/compiler"
	"golox/internal/debug"
	"golox/internal/opcode"
	"golox/internal/value"
	"os"
)

type InterpretResult = uint8

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
	InterpretNoResult
)

type VM struct {
	stack    []value.Value
	chunk    *chunk.Chunk
	ip       int
	stackTop int
	globals  map[string]value.Value
}

func NewVM() *VM {
	return &VM{
		stack: make([]value.Value, 0),
	}
}

func (vm *VM) push(value value.Value) InterpretResult {
	if vm.stackTop >= common.Uint24Max {
		vm.runtimeError("Stack overflow, tried to push with %v values on stack.", common.Uint24Max)
		return InterpretRuntimeError
	}

	vm.stack = append(vm.stack, value)
	vm.stackTop++
	return InterpretNoResult
}

func (vm *VM) pop() (value.Value, InterpretResult) {
	if vm.stackTop < 1 {
		vm.runtimeError("Stack underflow, tried to pop with no values on stack.")
		return value.NilVal{}, InterpretRuntimeError
	}
	vm.stackTop--
	poppedValue := vm.stack[vm.stackTop]
	vm.stack = vm.stack[:vm.stackTop]
	return poppedValue, InterpretNoResult
}

func (vm *VM) peek(distance int) value.Value {
	return vm.stack[vm.stackTop-1-distance]
}

func (vm *VM) Interpret(source *[]byte) InterpretResult {
	c := chunk.NewChunk()

	if !compiler.Compile(source, c) {
		c.Free()
		return InterpretCompileError
	}

	vm.chunk = c
	vm.ip = 0
	vm.globals = make(map[string]value.Value)

	result := vm.run()

	c.Free()
	return result
}

func (vm *VM) run() InterpretResult {
	for {
		if debug.TraceExecution {
			fmt.Printf("          ")
			for slot := 0; slot < vm.stackTop; slot++ {
				fmt.Printf("[ %s ]", vm.stack[slot].Stringify())
			}
			fmt.Printf("\n")
			debug.DisassembleInstruction(vm.chunk, vm.ip)
		}

		switch instruction := vm.readByte(); instruction {
		case opcode.Constant, opcode.ConstantLong:
			constant := vm.readConstant(instruction)
			pushResult := vm.push(constant)
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.Nil:
			pushResult := vm.push(value.NilVal{})
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.True:
			pushResult := vm.push(value.BoolVal(true))
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.False:
			pushResult := vm.push(value.BoolVal(false))
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.Pop:
			_, popResult := vm.pop()
			if popResult != InterpretNoResult {
				return popResult
			}
		case opcode.GetLocal, opcode.GetLocalLong:
			slot := vm.readIndex(instruction)
			pushResult := vm.push(vm.stack[slot])
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.SetLocal, opcode.SetLocalLong:
			slot := vm.readIndex(instruction)
			vm.stack[slot] = vm.peek(slot)
		case opcode.GetGlobal, opcode.GetGlobalLong:
			name := vm.readConstant(instruction).AsString()
			val, exists := vm.globals[name]
			_, popResult := vm.pop()
			if popResult != InterpretNoResult {
				return popResult
			}
			if !exists {
				vm.runtimeError("Undefined variable '%s'.", name)
				return InterpretRuntimeError
			}
			pushResult := vm.push(val)
			if pushResult != InterpretNoResult {
				return popResult
			}
		case opcode.DefineGlobal, opcode.DefineGlobalLong:
			name := vm.readConstant(instruction).AsString()
			val, popResult := vm.pop()
			if popResult != InterpretNoResult {
				return popResult
			}
			vm.globals[name] = val
		case opcode.SetGlobal, opcode.SetGlobalLong:
			name := vm.readConstant(instruction).AsString()
			_, exists := vm.globals[name]
			if exists {
				val, popResult := vm.pop()
				if popResult != InterpretNoResult {
					return popResult
				}
				vm.globals[name] = val
			} else {
				vm.runtimeError("Undefined variable '%s'.", name)
				return InterpretRuntimeError
			}
		case opcode.Equal:
			valB, popResultB := vm.pop()
			if popResultB != InterpretNoResult {
				return popResultB
			}
			valA, popResultA := vm.pop()
			if popResultA != InterpretNoResult {
				return popResultA
			}
			pushResult := vm.push(value.BoolVal((valA.IsEqual(valB))))
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.NotEqual:
			valB, popResultB := vm.pop()
			if popResultB != InterpretNoResult {
				return popResultB
			}
			valA, popResultA := vm.pop()
			if popResultA != InterpretNoResult {
				return popResultA
			}
			pushResult := vm.push(value.BoolVal((!valA.IsEqual(valB))))
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.Add:
			a := vm.peek(1)
			b := vm.peek(0)
			if a.IsString() && b.IsString() {
				valB, popResultB := vm.pop()
				if popResultB != InterpretNoResult {
					return popResultB
				}
				valA, popResultA := vm.pop()
				if popResultA != InterpretNoResult {
					return popResultA
				}
				pushResult := vm.push(value.StringVal(valA.AsString() + valB.AsString()))
				if pushResult != InterpretNoResult {
					return pushResult
				}
			} else if a.IsNumber() && b.IsNumber() {
				valB, popResultB := vm.pop()
				if popResultB != InterpretNoResult {
					return popResultB
				}
				valA, popResultA := vm.pop()
				if popResultA != InterpretNoResult {
					return popResultA
				}
				pushResult := vm.push(value.NumberVal(valA.AsNumber() + valB.AsNumber()))
				if pushResult != InterpretNoResult {
					return pushResult
				}
			} else {
				vm.runtimeError(
					"Operands must be two numbers or two strings.")
				return InterpretRuntimeError
			}
		case opcode.Greater, opcode.GreaterEqual, opcode.Less, opcode.LessEqual,
			opcode.Subtract, opcode.Multiply, opcode.Divide, opcode.Modulo:
			result := vm.binaryOP(instruction)
			if result != InterpretNoResult {
				return result
			}
		case opcode.Not:
			val, popResult := vm.pop()
			if popResult != InterpretNoResult {
				return popResult
			}
			pushResult := vm.push(value.BoolVal(val.IsFalsey()))
			if pushResult != InterpretNoResult {
				return pushResult
			}
		case opcode.Negate:
			if val := vm.peek(0); !val.IsNumber() {
				vm.runtimeError("Operand must be a number.")
				return InterpretRuntimeError
			} else {
				val, popResult := vm.pop()
				if popResult != InterpretNoResult {
					return popResult
				}
				pushResult := vm.push(value.NumberVal(-val.AsNumber()))
				if pushResult != InterpretNoResult {
					return pushResult
				}
			}
		case opcode.Print:
			val, popResult := vm.pop()
			if popResult != InterpretNoResult {
				return popResult
			}
			fmt.Printf("%s\n", val.Stringify())
		case opcode.Jump:
			offset := vm.readShort()
			vm.ip += offset
		case opcode.JumpIfFalse:
			offset := vm.readShort()
			if vm.peek(0).IsFalsey() {
				vm.ip += offset
			}
		case opcode.Loop:
			offset := vm.readShort()
			vm.ip -= offset
		case opcode.Return:
			return InterpretOk
		default:
			err := fmt.Sprintf("Unknown instruction %v", instruction)
			panic(err)
		}
	}
}

func (vm *VM) binaryOP(operator byte) InterpretResult {
	if !vm.peek(0).IsNumber() || !vm.peek(1).IsNumber() {
		vm.runtimeError("Operands must be numbers.")
		return InterpretRuntimeError
	}

	valB, popResultB := vm.pop()
	if popResultB != InterpretNoResult {
		return popResultB
	}
	valA, popResultA := vm.pop()
	if popResultA != InterpretNoResult {
		return popResultA
	}

	switch operator {
	case opcode.Greater:
		pushResult := vm.push(value.BoolVal(valA.AsNumber() > valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.GreaterEqual:
		pushResult := vm.push(value.BoolVal(valA.AsNumber() >= valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.Less:
		pushResult := vm.push(value.BoolVal(valA.AsNumber() < valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.LessEqual:
		pushResult := vm.push(value.BoolVal(valA.AsNumber() <= valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.Subtract:
		pushResult := vm.push(value.NumberVal(valA.AsNumber() - valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.Multiply:
		pushResult := vm.push(value.NumberVal(valA.AsNumber() * valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.Divide:
		pushResult := vm.push(value.NumberVal(valA.AsNumber() / valB.AsNumber()))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	case opcode.Modulo:
		pushResult := vm.push(value.NumberVal(float64(int(valA.AsNumber()) % int(valB.AsNumber()))))
		if pushResult != InterpretNoResult {
			return pushResult
		}
	default:
		err := fmt.Sprintf("Invalid binary operator %v", operator)
		panic(err)
	}

	return InterpretNoResult
}

func (vm *VM) readByte() byte {
	vm.ip++
	return vm.chunk.Code[vm.ip-1]
}

func (vm *VM) readShort() int {
	vm.ip += 2
	return int(uint16(vm.chunk.Code[vm.ip-2])<<8 | uint16(vm.chunk.Code[vm.ip-1]))
}

func (vm *VM) readConstant(op byte) value.Value {
	index := vm.readIndex(op)
	return vm.chunk.Constants[index]

}

func (vm *VM) readIndex(op byte) int {
	switch op {
	case opcode.Constant, opcode.DefineGlobal, opcode.GetGlobal,
		opcode.SetGlobal, opcode.GetLocal, opcode.SetLocal:
		return int(vm.readByte())
	case opcode.ConstantLong, opcode.DefineGlobalLong, opcode.GetGlobalLong,
		opcode.SetGlobalLong, opcode.GetLocalLong, opcode.SetLocalLong:
		index := uint32(vm.readByte()) << 16
		index |= uint32(vm.readByte()) << 8
		index |= uint32(vm.readByte())
		return int(index)
	default:
		msg := fmt.Sprintf("Invalid opcode '%v' for function vm.readIndex(op).", opcode.Name(op))
		panic(msg)
	}
}

func (vm *VM) runtimeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", vm.chunk.GetLine(vm.ip-1))
	vm.resetStack()
}

func (vm *VM) resetStack() {
	vm.stackTop = 0
}
