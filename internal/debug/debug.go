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
	case opcode.Constant, opcode.DefineGlobal, opcode.GetGlobal, opcode.SetGlobal:
		return constantInstruction(opcode.Name[op], c, offset)
	case opcode.ConstantLong, opcode.DefineGlobalLong, opcode.GetGlobalLong, opcode.SetGlobalLong:
		return constantLongInstruction(opcode.Name[op], c, offset)
	case opcode.Nil, opcode.True, opcode.False, opcode.Pop,
		opcode.Equal, opcode.NotEqual, opcode.Greater, opcode.GreaterEqual,
		opcode.Less, opcode.LessEqual, opcode.Add, opcode.Subtract,
		opcode.Multiply, opcode.Divide, opcode.Not, opcode.Modulo,
		opcode.Negate, opcode.Print, opcode.Return:
		return simpleInstruction(opcode.Name[op], offset)
	case opcode.GetLocal, opcode.SetLocal:
		return byteInstruction(opcode.Name[op], c, offset)
	case opcode.GetLocalLong, opcode.SetLocalLong:
		return byteInstructionLong(opcode.Name[op], c, offset)
	case opcode.Jump, opcode.JumpIfFalse:
		return jumpInstruction(opcode.Name[op], 1, c, offset)
	case opcode.Loop:
		return jumpInstruction(opcode.Name[op], -1, c, offset)
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
	constantIndex |= uint32(c.Code[offset+2]) << 8
	constantIndex |= uint32(c.Code[offset+3])
	fmt.Printf("%-16s %4d '%s'\n", name, constantIndex, c.Constants[constantIndex].Stringify())
	return offset + 4
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func byteInstruction(name string, c *chunk.Chunk, offset int) int {
	slot := c.Code[offset+1]
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 2
}

func byteInstructionLong(name string, c *chunk.Chunk, offset int) int {
	slot := uint32(c.Code[offset+1]) << 16
	slot |= uint32(c.Code[offset+2]) << 8
	slot |= uint32(c.Code[offset+3])
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 4
}

func jumpInstruction(name string, sign int, chunk *chunk.Chunk, offset int) int {
	jump := uint16(chunk.Code[offset+1]) << 8
	jump |= uint16(chunk.Code[offset+2])
	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*int(jump))
	return offset + 3
}
