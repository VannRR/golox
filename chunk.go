package main

import "fmt"

const (
	OP_CONSTANT uint8 = iota
	OP_CONSTANT_LONG
	OP_RETURN
)

const maxConstantIndex = 255
const maxConstantLongIndex = 16_777_215

type LineInfo struct {
	line  uint16
	count uint16
}

type Chunk struct {
	code      []uint8
	lineInfo  []LineInfo
	constants ValueArray
}

func NewChunk() *Chunk {
	return &Chunk{
		code:      make([]uint8, 0),
		lineInfo:  make([]LineInfo, 0),
		constants: NewValueArray(),
	}
}

func (c Chunk) Count() int {
	return len(c.code)
}

func (c Chunk) GetLine(codeIndex int) uint16 {
	if codeIndex < 0 || codeIndex >= len(c.code) {
		return 0
	}

	var cumulativeIndex uint16 = 0
	for _, l := range c.lineInfo {
		cumulativeIndex += l.count
		if cumulativeIndex > uint16(codeIndex) {
			return l.line
		}
	}

	return 0
}

func (c *Chunk) Free() {
	c.code = c.code[:0]
	c.lineInfo = c.lineInfo[:0]
	c.constants.Free()
}

func (c *Chunk) WriteConstant(value Value, line uint16) {
	if i := c.AddConstant(value); i <= maxConstantIndex {
		c.Write(OP_CONSTANT, line)
		c.Write(uint8(i), line)
	} else if i <= maxConstantLongIndex {
		c.Write(OP_CONSTANT_LONG, line)
		c.Write(uint8(i>>16), line)
		c.Write(uint8(i>>8), line)
		c.Write(uint8(i), line)
	} else {
		errMsg := fmt.Sprintf("Too many constants (%d), must be less than 16,777,216", i)
		panic(errMsg)
	}
}

func (c *Chunk) Write(byte uint8, line uint16) {
	c.code = append(c.code, byte)
	if last := len(c.lineInfo)-1; len(c.lineInfo) == 0 || c.lineInfo[last].line != line {
		c.lineInfo = append(c.lineInfo, LineInfo{line: line, count: 1})
	} else {
		c.lineInfo[last].count++
	}
}

func (c *Chunk) AddConstant(value Value) int {
	c.constants.Write(value)
	return c.constants.Count() - 1
}

func (c Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < c.Count(); {
		offset = disassembleInstruction(c, offset)
	}
}
