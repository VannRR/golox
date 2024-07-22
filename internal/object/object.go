package object

const (
	String ObjType = iota
)

type ObjType uint8

type Obj struct {
	Type ObjType
}

type ObjString struct {
	String string
	Obj
}

func (o *ObjString) Length() int {
	return len(o.String)
}
