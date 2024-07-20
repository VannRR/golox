package value

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	Bool ValueType = iota
	Nil
	Number
)

type ValueType = uint8

type Value struct {
	Type  ValueType
	bytes [8]byte
}

func NewBoolVal(v bool) Value {
	if v {
		return Value{
			Type:  Bool,
			bytes: [8]byte{1},
		}
	} else {
		return Value{
			Type:  Bool,
			bytes: [8]byte{0},
		}
	}

}

func NewNilVal() Value {
	return Value{
		Type: Nil,
	}
}

func NewNumberVal(v float64) Value {
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], math.Float64bits(v))
	return Value{
		Type:  Number,
		bytes: bytes,
	}
}

func (v *Value) IsBool() bool {
	return v.Type == Bool
}
func (v *Value) IsNil() bool {
	return v.Type == Nil
}
func (v *Value) IsNumber() bool {
	return v.Type == Number
}

func (v *Value) AsBool() bool {
	return v.bytes[0] == 1
}

func (v *Value) AsNumber() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(v.bytes[:]))
}

func (v Value) Print() {
	fmt.Printf("%g", v.AsNumber())
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
