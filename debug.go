package main

import "fmt"

func disassembleInstruction(c Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	
	if l := c.GetLine(offset); offset > 0 && l == c.GetLine(offset-1) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", l)
	}
	
	switch op := c.code[offset]; op {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_CONSTANT_LONG:
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", op)
		return offset + 1
	}
}

func constantInstruction(name string, c Chunk, offset int) int {
	constantIndex := c.code[offset+1]
	fmt.Printf("%-16s %4d '", name, constantIndex)
	c.constants[constantIndex].Print()
	fmt.Println("'")
	return offset + 2
}

func constantLongInstruction(name string, c Chunk, offset int) int {
	constantIndex := uint32(c.code[offset+1]) << 16
	constantIndex += uint32(c.code[offset+2]) << 8
	constantIndex += uint32(c.code[offset+3])
	fmt.Printf("%-16s %4d '", name, constantIndex)
	c.constants[constantIndex].Print()
	fmt.Println("'")
	return offset + 4
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
