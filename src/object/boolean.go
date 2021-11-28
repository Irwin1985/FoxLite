package object

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjType {
	return BooleanObj
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "True"
	}
	return "False"
}
