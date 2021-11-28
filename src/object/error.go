package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjType {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return e.Message
}
