package object

type Exit struct {
}

func (e *Exit) Type() ObjType {
	return ExitObj
}

func (e *Exit) Inspect() string {
	return "Exit"
}
