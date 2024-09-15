package chunk

import (
	"github.com/VannRR/golox/internal/opcode"
	"github.com/VannRR/golox/internal/value"
	"testing"
)

func Test_NewChunk(t *testing.T) {
	ch := NewChunk()

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 0)

	if len(ch.lineInfo) != 0 {
		t.Errorf("Expected LineInfo slice length of '0', got '%v'.", len(ch.lineInfo))
	}

}

func Test_GetLine(t *testing.T) {
	ch := NewChunk()

	ch.Write(0, 123)

	line := ch.GetLine(0)

	if line != 123 {
		t.Errorf("Expected line number '123', got '%v'.", line)
	}
}

func Test_Free(t *testing.T) {
	ch := NewChunk()

	index := ch.AddConstant(value.NilVal{})
	ch.WriteIndexWithCheck(index, opcode.Constant, 123)

	ch.Free()

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 0)

	if len(ch.lineInfo) != 0 {
		t.Errorf("Expected LineInfo slice length of '0', got '%v'.", len(ch.lineInfo))
	}
}

func Test_WriteIndexWithCheck(t *testing.T) {
	ch := NewChunk()
	globalVar := value.NumberVal(420)
	var line uint16 = 123
	index := ch.AddConstant(globalVar)
	ch.WriteIndexWithCheck(index, opcode.Constant, line)

	expectCodeCount(t, ch, 2)

	expectOpCodeAtIndex(t, ch, opcode.Constant, 1)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, globalVar, index)
}

func Test_AddConstant(t *testing.T) {
	ch := NewChunk()
	val := value.NumberVal(365)
	index := ch.AddConstant(val)

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, val, index)
}

func Test_Write(t *testing.T) {
	ch := NewChunk()
	op := opcode.Add
	ch.Write(op, 123)

	expectCodeCount(t, ch, 1)

	expectOpCodeAtIndex(t, ch, op, 0)
}

func expectCodeCount(t *testing.T, ch *Chunk, count int) {
	t.Helper()
	if ch.Count() != count {
		t.Fatalf("Expected byte code count of '%v', got '%v'.", count, ch.Count())
	}
}

func expectConstantCount(t *testing.T, ch *Chunk, count int) {
	t.Helper()
	if ch.Constants.Count() != count {
		t.Fatalf("Expected constants count of '%v', got '%v'.", count, ch.Constants.Count())
	}
}

func expectOpCodeAtIndex(t *testing.T, ch *Chunk, op byte, index int) {
	t.Helper()
	actual := ch.Code[index]
	if actual != op {
		name, exists := opcode.Name[actual]
		if !exists {
			t.Errorf("Unknown opcode '%v'.", actual)
		} else {
			t.Errorf("Expected opcode '%v', got '%v'.", opcode.Name[op], name)
		}
	}
}

func expectConstantAtIndex(t *testing.T, ch *Chunk, constant value.Value, index int) {
	t.Helper()
	if ch.Constants[index] != constant {
		t.Errorf("Expected constant '%v', got '%v'.", constant, ch.Constants[index])
	}
}
