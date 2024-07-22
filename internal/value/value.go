package value

import (
	"fmt"
	"golox/internal/obj-type"
	"golox/internal/val-type"
)

type Obj struct {
	Type objtype.ObjType
}

type ObjString struct {
	Obj
	String string
}

type Value struct {
	of   interface{}
	Type valtype.ValType
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
		Type: valtype.Obj,
		of: ObjString{
			Obj: Obj{
				Type: objtype.String,
			},
			String: s,
		},
	}
}

func (v *Value) IsFalsey() bool {
	switch v.Type {
	case valtype.Nil:
		return true
	case valtype.Bool:
		return v.of == false
	default:
		return false
	}
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

func (v *Value) IsObj() bool {
	return v.Type == valtype.Obj
}

func (v *Value) IsString() bool {
	if v.Type != valtype.Obj {
		return false
	}
	obj, ok := v.of.(ObjString)
	if !ok {
		return false
	}
	return obj.Obj.Type == objtype.String
}

func (v *Value) AsBool() bool {
	if b, ok := v.of.(bool); ok {
		return b
	}
	panic(fmt.Sprintf("expected bool, got '%T'", v.of))
}

func (v *Value) AsNumber() float64 {
	if n, ok := v.of.(float64); ok {
		return n
	}
	panic(fmt.Sprintf("expected float64, got '%T'", v.of))
}

func (v *Value) AsObj() Obj {
	if o, ok := v.of.(Obj); ok {
		return o
	}
	panic(fmt.Sprintf("expected Obj, got '%T'", v.of))
}

func (v *Value) ObjType() objtype.ObjType {
	switch obj := v.of.(type) {
	case ObjString:
		return objtype.String
	default:
		panic(fmt.Sprintf("unknown object type '%T'", obj))
	}
}

func (v *Value) AsString() ObjString {
	if s, ok := v.of.(ObjString); ok {
		return s
	}
	panic(fmt.Sprintf("expected ObjString, got '%T'", v.of))
}

func (v *Value) AsGoString() string {
	return v.AsString().String
}

func (v Value) Print() {
	switch v.Type {
	case valtype.Bool:
		fmt.Printf("%v", v.AsBool())
	case valtype.Nil:
		fmt.Print("nil")
	case valtype.Number:
		fmt.Printf("%v", v.AsNumber())
	case valtype.Obj:
		switch v.ObjType() {
		case objtype.String:
			fmt.Printf("%v", v.AsGoString())
		default:
			panic(fmt.Sprintf("cannot print unknown object type '%v'", v.ObjType()))
		}
	default:
		panic(fmt.Sprintf("cannot print unknown type '%v'", v.Type))
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
	case valtype.Obj:
		switch a.ObjType() {
		case objtype.String:
			return a.AsGoString() == b.AsGoString()
		default:
			panic(fmt.Sprintf("cannot compare unknown object type '%v' as equal", a.ObjType()))
		}
	default:
		panic(fmt.Sprintf("cannot compare unknown type '%v' as equal", a.Type))
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
