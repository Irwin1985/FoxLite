package object

import "fmt"

type Integer struct {
	Value float64
}

func (i *Integer) Type() ObjType {
	return IntegerObj
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%v", i.Value)
}
