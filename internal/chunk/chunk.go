package chunk

import (
	"fmt"
	"golox/internal/common"
	"golox/internal/value"
)

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

func (c *Chunk) WriteIndexWithCheck(index int, opcode byte, line uint16) {
	if index <= common.Uint8Max {
		c.Write(opcode, line)
		c.Write(byte(index), line)
	} else if index <= common.Uint24Max {
		c.Write(opcode+1, line)
		c.Write(byte(index>>16), line)
		c.Write(byte(index>>8), line)
		c.Write(byte(index), line)
	} else {
		msg := fmt.Sprintf("Too many constants (%d), must be less than (%d)", index, common.Uint24Max+1)
		panic(msg)
	}
}

func (c *Chunk) AddConstant(value value.Value) int {
	c.Constants.Write(value)
	return c.Constants.Count() - 1
}

func (c *Chunk) Write(byte byte, line uint16) {
	c.Code = append(c.Code, byte)
	if last := len(c.lineInfo) - 1; len(c.lineInfo) == 0 || c.lineInfo[last].line != line {
		c.lineInfo = append(c.lineInfo, LineInfo{line: line, count: 1})
	} else {
		c.lineInfo[last].count++
	}
}
