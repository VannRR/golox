package value

import (
	"fmt"
)

type Value interface {
	fmt.Stringer
	IsEqual(Value) bool
	IsFalsey() bool
	IsType(Value) bool
	IsBool() bool
	IsNil() bool
	IsNumber() bool
	IsString() bool
	IsFunction() bool
}

type NilVal struct{}

func (n NilVal) String() string { return "nil" }

func (n NilVal) IsEqual(other Value) bool { return other.IsNil() }

func (n NilVal) IsFalsey() bool { return true }

func (n NilVal) IsType(other Value) bool { return other.IsNil() }
func (n NilVal) IsNil() bool             { return true }
func (n NilVal) IsBool() bool            { return false }
func (n NilVal) IsNumber() bool          { return false }
func (n NilVal) IsString() bool          { return false }
func (n NilVal) IsFunction() bool        { return false }

type BoolVal bool

func (b BoolVal) String() string { return fmt.Sprint(bool(b)) }

func (b BoolVal) IsEqual(other Value) bool {
	if !other.IsBool() {
		return false
	}
	return b == other.(BoolVal)
}

func (b BoolVal) IsFalsey() bool { return !bool(b) }

func (b BoolVal) IsType(other Value) bool { return other.IsBool() }
func (b BoolVal) IsBool() bool            { return true }
func (b BoolVal) IsNil() bool             { return false }
func (b BoolVal) IsNumber() bool          { return false }
func (b BoolVal) IsString() bool          { return false }
func (b BoolVal) IsFunction() bool        { return false }

type NumberVal float64

func (n NumberVal) String() string { return fmt.Sprint(float64(n)) }

func (n NumberVal) IsEqual(other Value) bool {
	if !other.IsNumber() {
		return false
	}
	return n == other.(NumberVal)
}

func (n NumberVal) IsFalsey() bool { return false }

func (n NumberVal) IsType(other Value) bool { return other.IsNumber() }
func (n NumberVal) IsBool() bool            { return false }
func (n NumberVal) IsNil() bool             { return false }
func (n NumberVal) IsNumber() bool          { return true }
func (n NumberVal) IsString() bool          { return false }
func (n NumberVal) IsFunction() bool        { return false }

type ValueArray []Value

func NewValueArray() ValueArray {
	return make(ValueArray, 0)
}

func (v ValueArray) Count() int {
	return len(v)
}

func (v *ValueArray) Free() {
	*v = (*v)[:0]
}

func (v *ValueArray) Write(value Value) {
	*v = append(*v, value)
}
