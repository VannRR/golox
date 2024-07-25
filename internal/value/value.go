package value

import (
	"fmt"
)

type Value interface {
	IsEqual(Value) bool
	IsFalsey() bool
	IsBool() bool
	IsNil() bool
	IsNumber() bool
	IsString() bool
	AsBool() bool
	AsNumber() float64
	AsString() string
	Stringify() string
}

type NilVal struct{}

func (a NilVal) IsEqual(b Value) bool {
	return b.IsNil()
}

func (n NilVal) IsFalsey() bool {
	return true
}

func (n NilVal) IsNil() bool {
	return true
}

func (n NilVal) IsBool() bool {
	return false
}

func (n NilVal) IsNumber() bool {
	return false
}

func (n NilVal) IsString() bool {
	return false
}

func (n NilVal) AsBool() bool {
	panic("Expected BoolVal, got NilVal")
}

func (n NilVal) AsNumber() float64 {
	panic("Expected NumberVal, got NilVal")
}

func (n NilVal) AsString() string {
	panic("Expected StringVal, got NilVal")
}

func (n NilVal) Stringify() string {
	return "nil"
}

type BoolVal bool

func (a BoolVal) IsEqual(b Value) bool {
	if !b.IsBool() {
		return false
	}
	return a.AsBool() == b.AsBool()
}

func (b BoolVal) IsFalsey() bool {
	return !bool(b)
}

func (b BoolVal) IsBool() bool {
	return true
}

func (b BoolVal) IsNil() bool {
	return false
}

func (b BoolVal) IsNumber() bool {
	return false
}

func (b BoolVal) IsString() bool {
	return false
}

func (b BoolVal) AsBool() bool {
	return bool(b)
}

func (b BoolVal) AsNumber() float64 {
	panic("Expected NumberVal, got BoolVal")
}

func (b BoolVal) AsString() string {
	panic("Expected StringVal, got BoolVal")
}

func (b BoolVal) Stringify() string {
	if b.AsBool() {
		return "true"
	} else {
		return "false"
	}
}

type NumberVal float64

func (a NumberVal) IsEqual(b Value) bool {
	if !b.IsNumber() {
		return false
	}
	return a.AsNumber() == b.AsNumber()
}

func (n NumberVal) IsFalsey() bool {
	return false
}

func (n NumberVal) IsBool() bool {
	return false
}

func (n NumberVal) IsNil() bool {
	return false
}

func (n NumberVal) IsNumber() bool {
	return true
}

func (n NumberVal) IsString() bool {
	return false
}

func (n NumberVal) AsBool() bool {
	panic("Expected BoolVal, got NumberVal")
}

func (n NumberVal) AsNumber() float64 {
	return float64(n)
}

func (n NumberVal) AsString() string {
	panic("Expected StringVal, got NumberVal")
}

func (n NumberVal) Stringify() string {
	return fmt.Sprintf("%v", n.AsNumber())
}

type StringVal string

func (a StringVal) IsEqual(b Value) bool {
	if !b.IsString() {
		return false
	}
	return a.AsString() == b.AsString()
}

func (s StringVal) IsFalsey() bool {
	return false
}

func (s StringVal) IsBool() bool {
	return false
}

func (s StringVal) IsNil() bool {
	return false
}

func (s StringVal) IsNumber() bool {
	return false
}

func (s StringVal) IsString() bool {
	return true
}

func (s StringVal) AsBool() bool {
	panic("Expected BoolVal, got StringVal")
}

func (s StringVal) AsNumber() float64 {
	panic("Expected NumberVal, got StringVal")
}

func (s StringVal) AsString() string {
	return string(s)
}

func (s StringVal) Stringify() string {
	return s.AsString()
}

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
