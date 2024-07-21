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
	case opcode.Constant:
		return constantInstruction("Constant", c, offset)
	case opcode.ConstantLong:
		return constantLongInstruction("ConstantLong", c, offset)
	case opcode.Nil:
		return simpleInstruction("Nil", offset)
	case opcode.True:
		return simpleInstruction("True", offset)
	case opcode.False:
		return simpleInstruction("False", offset)
	case opcode.Add:
		return simpleInstruction("Add", offset)
	case opcode.Subtract:
		return simpleInstruction("Subtract", offset)
	case opcode.Multiply:
		return simpleInstruction("Multiply", offset)
	case opcode.Divide:
		return simpleInstruction("Divide", offset)
	case opcode.Not:
		return simpleInstruction("Not", offset)
	case opcode.Modulo:
		return simpleInstruction("Modulo", offset)
	case opcode.Negate:
		return simpleInstruction("Negate", offset)
	case opcode.Return:
		return simpleInstruction("Return", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", op)
		return offset + 1
	}
}

func constantInstruction(name string, c *chunk.Chunk, offset int) int {
	constantIndex := c.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constantIndex)
	c.Constants[constantIndex].Print()
	fmt.Println("'")
	return offset + 2
}

func constantLongInstruction(name string, c *chunk.Chunk, offset int) int {
	constantIndex := uint32(c.Code[offset+1]) << 16
	constantIndex += uint32(c.Code[offset+2]) << 8
	constantIndex += uint32(c.Code[offset+3])
	fmt.Printf("%-16s %4d '", name, constantIndex)
	c.Constants[constantIndex].Print()
	fmt.Println("'")
	return offset + 4
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
