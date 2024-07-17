package chunk

import (
	"fmt"
	"golox/internal/value"
)

const (
	OP_CONSTANT byte = iota
	OP_CONSTANT_LONG
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

const maxConstantIndex = 255
const maxConstantLongIndex = 16_777_215

type LineInfo struct {
	line  uint16
	count uint16
}

type Chunk struct {
	Code      []byte
	lineInfo  []LineInfo
	Constants value.ValueArray
}

func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0),
		lineInfo:  make([]LineInfo, 0),
		Constants: value.NewValueArray(),
	}
}

func (c *Chunk) Count() int {
	return len(c.Code)
}

func (c *Chunk) GetLine(codeIndex int) uint16 {
	if codeIndex < 0 || codeIndex >= len(c.Code) {
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
	c.Code = c.Code[:0]
	c.lineInfo = c.lineInfo[:0]
	c.Constants.Free()
}

func (c *Chunk) WriteConstant(value value.Value, line uint16) {
	if i := c.AddConstant(value); i <= maxConstantIndex {
		c.Write(OP_CONSTANT, line)
		c.Write(byte(i), line)
	} else if i <= maxConstantLongIndex {
		c.Write(OP_CONSTANT_LONG, line)
		c.Write(byte(i>>16), line)
		c.Write(byte(i>>8), line)
		c.Write(byte(i), line)
	} else {
		errMsg := fmt.Sprintf("Too many constants (%d), must be less than 16,777,216", i)
		panic(errMsg)
	}
}

func (c *Chunk) Write(byte byte, line uint16) {
	c.Code = append(c.Code, byte)
	if last := len(c.lineInfo) - 1; len(c.lineInfo) == 0 || c.lineInfo[last].line != line {
		c.lineInfo = append(c.lineInfo, LineInfo{line: line, count: 1})
	} else {
		c.lineInfo[last].count++
	}
}

func (c *Chunk) AddConstant(value value.Value) int {
	c.Constants.Write(value)
	return c.Constants.Count() - 1
}
