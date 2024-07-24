package chunk_test

import (
	"golox/internal/chunk"
	"golox/internal/opcode"
	"golox/internal/value"
	"testing"
)

func Test_NewChunk(t *testing.T) {
	ch := chunk.NewChunk()

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 0)

	if len(ch.LineInfo) != 0 {
		t.Errorf("Expected LineInfo slice length of 0, got %v.", len(ch.LineInfo))
	}

}

func Test_GetLine(t *testing.T) {
	ch := chunk.NewChunk()

	ch.Write(0, 123)

	line := ch.GetLine(0)

	if line != 123 {
		t.Errorf("Expected line number 123, got %v.", line)
	}
}

func Test_Free(t *testing.T) {
	ch := chunk.NewChunk()

	ch.WriteConstant(value.NilVal{}, 123)

	ch.Free()

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 0)

	if len(ch.LineInfo) != 0 {
		t.Errorf("Expected LineInfo slice length of 0, got %v.", len(ch.LineInfo))
	}
}

func Test_WriteDefineGlobalVar(t *testing.T) {
	ch := chunk.NewChunk()
	globalVar := value.NumberVal(420)
	var line uint16 = 123
	index := ch.WriteConstant(globalVar, line)
	ch.WriteDefineGlobalVar(index, line)

	expectCodeCount(t, ch, 4)

	expectOpCodeAtIndex(t, ch, opcode.DefineGlobal, 2)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, globalVar, index)
}

func Test_WriteGetGlobalVar(t *testing.T) {
	ch := chunk.NewChunk()
	globalVar := value.NumberVal(21)
	var line uint16 = 123
	index := ch.WriteConstant(globalVar, line)
	ch.WriteGetGlobalVar(index, line)

	expectCodeCount(t, ch, 4)

	expectOpCodeAtIndex(t, ch, opcode.GetGlobal, 2)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, globalVar, index)
}

func Test_WriteConstant(t *testing.T) {
	ch := chunk.NewChunk()
	val := value.NumberVal(21)
	index := ch.WriteConstant(val, 123)

	expectCodeCount(t, ch, 2)

	expectOpCodeAtIndex(t, ch, opcode.Constant, 0)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, val, index)
}

func Test_AddConstant(t *testing.T) {
	ch := chunk.NewChunk()
	val := value.NumberVal(365)
	index := ch.AddConstant(val)

	expectCodeCount(t, ch, 0)

	expectConstantCount(t, ch, 1)

	expectConstantAtIndex(t, ch, val, index)
}

func Test_Write(t *testing.T) {
	ch := chunk.NewChunk()
	op := opcode.Add
	ch.Write(op, 123)

	expectCodeCount(t, ch, 1)

	expectOpCodeAtIndex(t, ch, op, 0)
}

func expectCodeCount(t *testing.T, ch *chunk.Chunk, count int) {
	if ch.Count() != count {
		t.Errorf("Expected byte code count of %v, got %v.", count, ch.Count())
	}
}

func expectConstantCount(t *testing.T, ch *chunk.Chunk, count int) {
	if ch.Constants.Count() != count {
		t.Errorf("Expected constants count of %v, got %v.", count, ch.Constants.Count())
	}
}

func expectOpCodeAtIndex(t *testing.T, ch *chunk.Chunk, op byte, index int) {
	if ch.Code[index] != op {
		t.Errorf("Expected opcode '%v', got '%v'.", opcode.Name(op), opcode.Name(ch.Code[index]))
	}
}

func expectConstantAtIndex(t *testing.T, ch *chunk.Chunk, constant value.Value, index int) {
	if ch.Constants[index] != constant {
		t.Errorf("Expected constant '%v', got '%v'.", constant, ch.Constants[index])
	}
}
