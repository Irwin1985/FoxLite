package object

type Return struct {
	Value Object
}

func (r *Return) Type() ObjType {
	return ReturnObj
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}
