package debug

import (
	"fmt"
	"golox/internal/chunk"
)

func DisassembleInstruction(c chunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if l := c.GetLine(offset); offset > 0 && l == c.GetLine(offset-1) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", l)
	}

	switch op := c.Code[offset]; op {
	case chunk.OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case chunk.OP_CONSTANT_LONG:
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)
	case chunk.OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case chunk.OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case chunk.OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case chunk.OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case chunk.OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case chunk.OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", op)
		return offset + 1
	}
}

func constantInstruction(name string, c chunk.Chunk, offset int) int {
	constantIndex := c.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constantIndex)
	c.Constants[constantIndex].Print()
	fmt.Println("'")
	return offset + 2
}

func constantLongInstruction(name string, c chunk.Chunk, offset int) int {
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
