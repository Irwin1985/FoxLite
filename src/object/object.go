package object

type ObjType int

const (
	IntegerObj ObjType = iota
	StringObj
	BooleanObj
	NullObj
	FuncObj
	ReturnObj
	ErrorObj
)

type Object interface {
	Type() ObjType
	Inspect() string
}

func NewError(msg string) *Error {
	return &Error{Message: msg}
}

func TypeToStr(t ObjType) string {
	switch t {
	case IntegerObj:
		return "int Literal"
	case StringObj:
		return "string"
	case BooleanObj:
		return "bool"
	case NullObj:
		return "null"
	default:
		return ""
	}
}
