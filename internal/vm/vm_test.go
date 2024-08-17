package vm

import (
	"golox/internal/chunk"
	"golox/internal/object"
	"golox/internal/opcode"
	"golox/internal/value"
	"testing"
)

func Test_NewVM(t *testing.T) {
	vm := NewVM()
	if vm == nil {
		t.Errorf("Expected a non-nil VM, got nil")
	}
}

func Test_push(t *testing.T) {
	vm := NewVM()
	val := value.NumberVal(42)

	pushResult := vm.push(val)
	if pushResult != InterpretNoResult {
		t.Fatalf("Expected no overflow with push")
	}

	if vm.stackTop != 1 {
		t.Errorf("Expected stackTop to be 1, got %d", vm.stackTop)
	}

	top := vm.stack[vm.stackTop-1]
	if top != val {
		t.Errorf("Expected pushed value to be %v, got %v", val, top)
	}
}

func Test_push_overflow(t *testing.T) {
	vm := NewVM()
	val := value.NilVal{}

	vm.stackTop = 16_777_215

	pushResult := vm.push(val)

	if pushResult != InterpretRuntimeError {
		t.Errorf("Expected overflow runtime error with push")
	}
}

func Test_pop(t *testing.T) {
	vm := NewVM()
	val := value.NumberVal(42)

	pushResult := vm.push(val)
	if pushResult != InterpretNoResult {
		t.Fatalf("Expected no overflow with push")
	}

	popped, popResult := vm.pop()
	if popResult != InterpretNoResult {
		t.Fatalf("Expected no underflow with pop")
	}

	if vm.stackTop != 0 {
		t.Errorf("Expected stackTop to be 0 after pop, got %d", vm.stackTop)
	}
	if popped != val {
		t.Errorf("Expected popped value to be %v, got %v", val, popped)
	}
}

func Test_pop_underflow(t *testing.T) {
	vm := NewVM()

	_, popResult := vm.pop()

	if popResult != InterpretRuntimeError {
		t.Errorf("Expected underflow runtime error with pop")
	}
}

func Test_peek(t *testing.T) {
	vm := NewVM()
	val := value.NumberVal(42)

	pushResult := vm.push(val)
	if pushResult != InterpretNoResult {
		t.Fatalf("Expected no overflow with push")
	}

	peeked := vm.peek(0)

	if peeked != val {
		t.Errorf("Expected peeked value to be %v, got %v", val, peeked)
	}
}

func Test_Interpret(t *testing.T) {
	vm := NewVM()
	source := []byte("print 42;")
	result := vm.Interpret(&source)

	if result != InterpretOk {
		t.Errorf("Expected Interpret result to be InterpretOk, got %d", result)
	}
}

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected InterpretResult
	}{
		{
			name:     "constant",
			source:   "123;",
			expected: InterpretOk,
		},
		{
			name:     "nil",
			source:   "nil;",
			expected: InterpretOk,
		},
		{
			name:     "true",
			source:   "true;",
			expected: InterpretOk,
		},
		{
			name:     "false",
			source:   "false;",
			expected: InterpretOk,
		},
		{
			name:     "define and get local",
			source:   "{var local = 123; local;}",
			expected: InterpretOk,
		},
		{
			name:     "define and set local",
			source:   "{var local = 123; local = 321;}",
			expected: InterpretOk,
		},

		{
			name:     "define and get global",
			source:   "var global = 123; global;",
			expected: InterpretOk,
		},
		{
			name:     "define and set global",
			source:   "var global = 123; global = 321;",
			expected: InterpretOk,
		},

		{
			name:     "equal",
			source:   "123 == 123;",
			expected: InterpretOk,
		},
		{
			name:     "not equal",
			source:   "123 != 456;",
			expected: InterpretOk,
		},
		{
			name:     "add",
			source:   "123 + 456;",
			expected: InterpretOk,
		},
		{
			name:     "greater",
			source:   "123 > 456;",
			expected: InterpretOk,
		},
		{
			name:     "less",
			source:   "123 < 456;",
			expected: InterpretOk,
		},
		{
			name:     "not",
			source:   "!false;",
			expected: InterpretOk,
		},
		{
			name:     "negate",
			source:   "-123;",
			expected: InterpretOk,
		},
		{
			name:     "print",
			source:   "print 123;",
			expected: InterpretOk,
		},
		{
			name:     "if true",
			source:   "if (true) print 1;",
			expected: InterpretOk,
		},
		{
			name:     "if false else",
			source:   "if (false) print 1; else print 2;",
			expected: InterpretOk,
		},
		{
			name:     "for i < 3 loop",
			source:   "for (var i = 0; i < 3; i = i + 1) print i;",
			expected: InterpretOk,
		},
		{
			name:     "while false loop",
			source:   "while (false) print 1;",
			expected: InterpretOk,
		},
		{
			name:     "while i < 3 loop",
			source:   "var i = 0; while (i < 3) i = i + 1; print i;",
			expected: InterpretOk,
		},
		{
			name:     "operands must be two numbers or two strings",
			source:   "123 + true;",
			expected: InterpretRuntimeError,
		},
		{
			name:     "operand must be a number",
			source:   `"cow" - 123;`,
			expected: InterpretRuntimeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.source)
			vm := NewVM()
			result := vm.Interpret(&source)
			if result != tt.expected {
				t.Errorf("Expected Interpret result to be %d, got %d", tt.expected, result)
			}
		})
	}
}

func Test_readByte(t *testing.T) {
	vm := &VM{
		ip: 0,
		chunk: &chunk.Chunk{
			Code: []byte{0x01, 0x02, 0x03},
		},
	}

	expected := byte(0x01)
	actual := vm.readByte()

	if actual != expected {
		t.Errorf("Expected %x, but got %x", expected, actual)
	}
}

func Test_readShort(t *testing.T) {
	vm := &VM{
		ip: 0,
		chunk: &chunk.Chunk{
			Code: []byte{0x1, 0x2},
		},
	}

	expected := 0x102
	actual := vm.readShort()

	if actual != expected {
		t.Errorf("Expected %x, but got %x", expected, actual)
	}
}

func Test_readConstant(t *testing.T) {
	vm := &VM{
		ip: 0,
		chunk: &chunk.Chunk{
			Code: []byte{opcode.Constant, 0},
			Constants: []value.Value{
				value.NumberVal(42),
			},
		},
	}

	expected := value.NumberVal(42)
	actual := vm.readConstant(opcode.Constant)

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func Test_readIndex(t *testing.T) {
	vm := &VM{
		ip: 0,
		chunk: &chunk.Chunk{
			Code:      []byte{21},
			Constants: []value.Value{},
		},
	}

	expected := 21
	actual := vm.readIndex(opcode.GetLocal)

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func Test_readIndex_Long(t *testing.T) {

	expected := 22222

	vm := &VM{
		ip: 0,
		chunk: &chunk.Chunk{
			Code: []byte{
				byte(expected >> 16),
				byte(expected >> 8),
				byte(expected),
			},
			Constants: []value.Value{},
		},
	}

	actual := vm.readIndex(opcode.GetLocalLong)

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func Test_add_numbers(t *testing.T) {
	vm := NewVM()
	a := 1
	b := 3
	vm.push(value.NumberVal(a))
	vm.push(value.NumberVal(b))

	result := vm.add()

	if result != InterpretNoResult {
		t.Errorf("Expected InterpretNoResult, got %v", result)
	}

	expected := value.NumberVal(a + b)

	actual, popResult := vm.pop()
	if popResult != InterpretNoResult {
		t.Fatalf("Expected no underflow with pop")
	}

	if actual != expected {
		t.Errorf("Expected (%v + %v) == %v, got %v", a, b, expected, actual)
	}
}

func Test_add_strings(t *testing.T) {
	vm := NewVM()
	a := "foo"
	b := "bar"
	vm.push(object.ObjString(a))
	vm.push(object.ObjString(b))

	result := vm.add()

	if result != InterpretNoResult {
		t.Errorf("Expected InterpretNoResult, got %v", result)
	}

	expected := object.ObjString(a + b)

	actual, popResult := vm.pop()
	if popResult != InterpretNoResult {
		t.Fatalf("Expected no underflow with pop")
	}

	if actual != expected {
		t.Errorf("Expected (%v + %v) == %v, got %v", a, b, expected, actual)
	}
}

func Test_add_error(t *testing.T) {
	vm := NewVM()
	a := "foo"
	b := 1
	vm.push(object.ObjString(a))
	vm.push(value.NumberVal(b))

	result := vm.add()

	if result != InterpretRuntimeError {
		t.Errorf("Expected InterpretRuntimeError, got %v", result)
	}
}

func Test_binaryOP(t *testing.T) {
	checkBinaryOp(t, 10, 4, opcode.Greater, value.BoolVal(true))
	checkBinaryOp(t, 10, 10, opcode.GreaterEqual, value.BoolVal(true))
	checkBinaryOp(t, 10, 4, opcode.GreaterEqual, value.BoolVal(true))
	checkBinaryOp(t, 4, 10, opcode.Less, value.BoolVal(true))
	checkBinaryOp(t, 10, 10, opcode.LessEqual, value.BoolVal(true))
	checkBinaryOp(t, 4, 10, opcode.LessEqual, value.BoolVal(true))
	checkBinaryOp(t, 10, 4, opcode.Subtract, value.NumberVal(6))
	checkBinaryOp(t, 10, 4, opcode.Multiply, value.NumberVal(40))
	checkBinaryOp(t, 10, 4, opcode.Modulo, value.NumberVal(2))
}

func Test_resetStack(t *testing.T) {
	vm := NewVM()
	vm.push(value.NilVal{})
	vm.resetStack()
	if vm.stackTop != 0 {
		t.Errorf("Expected stackTop to = 0, got %v", vm.stackTop)
	}
}

func checkBinaryOp(t *testing.T, a int, b int, operation byte, expected value.Value) {
	t.Helper()
	vm := NewVM()
	vm.push(value.NumberVal(a))
	vm.push(value.NumberVal(b))

	vm.binaryOP(operation)

	actual, popResult := vm.pop()
	if popResult != InterpretNoResult {
		t.Fatalf("Expected no underflow with pop")
	}

	if actual != expected {
		t.Errorf("Expected (%v %v %v) == %v, got %v", a, opcode.Name[operation], b, expected, actual)
	}
}
