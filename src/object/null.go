package object

type Null struct {
}

func (n *Null) Type() ObjType {
	return NullObj
}

func (n *Null) Inspect() string {
	return "Null"
}
