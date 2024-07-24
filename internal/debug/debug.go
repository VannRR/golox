package debug

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/opcode"
)

const PrintCode bool = true
const TraceExecution bool = true

func DisassembleChunk(c *chunk.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < c.Count(); {
		offset = DisassembleInstruction(c, offset)
	}
}

func DisassembleInstruction(c *chunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if l := c.GetLine(offset); offset > 0 && l == c.GetLine(offset-1) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", l)
	}

	switch op := c.Code[offset]; op {
	case opcode.Constant, opcode.DefineGlobal, opcode.GetGlobal:
		return constantInstruction(opcode.Name(op), c, offset)
	case opcode.ConstantLong, opcode.DefineGlobalLong, opcode.GetGlobalLong:
		return constantLongInstruction(opcode.Name(op), c, offset)
	case opcode.Nil, opcode.True, opcode.False, opcode.Pop,
		opcode.Equal, opcode.NotEqual, opcode.Greater, opcode.GreaterEqual,
		opcode.Less, opcode.LessEqual, opcode.Add, opcode.Subtract,
		opcode.Multiply, opcode.Divide, opcode.Not, opcode.Modulo,
		opcode.Negate, opcode.Print, opcode.Return:
		return simpleInstruction(opcode.Name(op), offset)
	default:
		fmt.Printf("Unknown opcode %d\n", op)
		return offset + 1
	}
}

func constantInstruction(name string, c *chunk.Chunk, offset int) int {
	constantIndex := c.Code[offset+1]
	fmt.Printf("%-16s %4d '%s'\n", name, constantIndex, c.Constants[constantIndex].Stringify())
	return offset + 2
}

func constantLongInstruction(name string, c *chunk.Chunk, offset int) int {
	constantIndex := uint32(c.Code[offset+1]) << 16
	constantIndex += uint32(c.Code[offset+2]) << 8
	constantIndex += uint32(c.Code[offset+3])
	fmt.Printf("%-16s %4d '%s'\n", name, constantIndex, c.Constants[constantIndex].Stringify())
	return offset + 4
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
