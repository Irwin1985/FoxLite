package object

type Loop struct {
}

func (l *Loop) Type() ObjType {
	return LoopObj
}

func (l *Loop) Inspect() string {
	return "Loop"
}
