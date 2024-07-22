package value

import (
	"fmt"
	"golox/internal/val-type"
)

type Value struct {
	of   interface{}
	Type valtype.Type
}

func NewBool(v bool) Value {
	return Value{
		Type: valtype.Bool,
		of:   v,
	}
}

func NewNil() Value {
	return Value{
		Type: valtype.Nil,
	}
}

func NewNumber(v float64) Value {
	return Value{
		Type: valtype.Number,
		of:   v,
	}
}

func NewObjString(s string) Value {
	return Value{
		Type: valtype.ObjString,
		of:   s,
	}
}

func (v *Value) IsFalsey() bool {
	return v.Type == valtype.Nil || (v.Type == valtype.Bool && v.of == false)
}

func (v *Value) IsBool() bool {
	return v.Type == valtype.Bool
}

func (v *Value) IsNil() bool {
	return v.Type == valtype.Nil
}

func (v *Value) IsNumber() bool {
	return v.Type == valtype.Number
}

func (v *Value) IsString() bool {
	return v.Type == valtype.ObjString
}

func (v *Value) AsBool() bool {
	return v.of.(bool)
}

func (v *Value) AsNumber() float64 {
	return v.of.(float64)
}

func (v *Value) AsString() string {
	return v.of.(string)
}

func (v Value) Print() {
	switch v.Type {
	case valtype.Bool:
		fmt.Printf("%v", v.AsBool())
	case valtype.Nil:
		fmt.Print("nil")
	case valtype.Number:
		fmt.Printf("%v", v.AsNumber())
	case valtype.ObjString:
		fmt.Printf("%v", v.AsString())
	default:
		msg := fmt.Sprintf("cannot print unknown type '%v'", v.Type)
		panic(msg)
	}
}

func (a *Value) IsEqual(b *Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case valtype.Bool:
		return a.AsBool() == b.AsBool()
	case valtype.Nil:
		return true
	case valtype.Number:
		return a.AsNumber() == b.AsNumber()
	case valtype.ObjString:
		return a.AsString() == b.AsString()
	default:
		msg := fmt.Sprintf("cannot compare unknown type '%v' as equal", a.Type)
		panic(msg)
	}
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
