package object

type None struct {
}

func (n *None) Type() ObjType {
	return NoneObj
}

func (n *None) Inspect() string {
	return ""
}
