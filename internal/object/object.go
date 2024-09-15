package object

import (
	"fmt"
	"github.com/VannRR/golox/internal/chunk"
	"github.com/VannRR/golox/internal/value"
)

type ObjString string

func (s ObjString) String() string { return string(s) }

func (s ObjString) IsEqual(other value.Value) bool {
	if !other.IsString() {
		return false
	}
	return s == other.(ObjString)
}

func (s ObjString) IsFalsey() bool { return false }

func (s ObjString) IsType(other value.Value) bool { return other.IsString() }
func (s ObjString) IsBool() bool                  { return false }
func (s ObjString) IsNil() bool                   { return false }
func (s ObjString) IsNumber() bool                { return false }
func (s ObjString) IsString() bool                { return true }
func (s ObjString) IsFunction() bool              { return false }

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

func (f ObjFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.name)
}

func (f ObjFunction) IsEqual(other value.Value) bool { return false }

func (f ObjFunction) IsFalsey() bool { return false }

func (f ObjFunction) IsType(other value.Value) bool { return other.IsFunction() }
func (f ObjFunction) IsBool() bool                  { return false }
func (f ObjFunction) IsNil() bool                   { return false }
func (f ObjFunction) IsNumber() bool                { return false }
func (f ObjFunction) IsString() bool                { return false }
func (f ObjFunction) IsFunction() bool              { return true }

func (f *ObjFunction) Free() { f.chunk.Free() }
