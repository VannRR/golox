package object

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/value"
)

type ObjString string

func (s ObjString) IsType(other value.Value) bool {
	return other.IsString()
}

func (s ObjString) IsEqual(other value.Value) bool {
	if !other.IsString() {
		return false
	}
	return s == other.(ObjString)
}

func (s ObjString) IsFalsey() bool {
	return false
}

func (s ObjString) IsBool() bool {
	return false
}

func (s ObjString) IsNil() bool {
	return false
}

func (s ObjString) IsNumber() bool {
	return false
}

func (s ObjString) IsString() bool {
	return true
}

func (s ObjString) IsFunction() bool {
	return false
}

func (s ObjString) Stringify() string {
	return string(s)
}

type ObjFunction struct {
	arity int
	chunk chunk.Chunk
	name  string
}

func NewFunction() *ObjFunction {
	return &ObjFunction{
		arity: 0,
		name:  "",
		chunk: *chunk.NewChunk(),
	}
}

func (f ObjFunction) IsType(other value.Value) bool {
	return other.IsFunction()
}

func (f ObjFunction) IsEqual(other value.Value) bool {
	if !other.IsFunction() {
		return false
	}
	return fmt.Sprint(f) == fmt.Sprint(other.(ObjFunction))
}

func (f ObjFunction) IsFalsey() bool {
	return false
}

func (f ObjFunction) IsBool() bool {
	return false
}

func (f ObjFunction) IsNil() bool {
	return false
}

func (f ObjFunction) IsNumber() bool {
	return false
}

func (f ObjFunction) IsString() bool {
	return false
}

func (f ObjFunction) IsFunction() bool {
	return true
}

func (f ObjFunction) Stringify() string {
	return fmt.Sprintf("<fn %s>", f.name)
}
