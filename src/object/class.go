package object

type Class struct {
	Name       string
	Properties map[string]Object
	Methods    map[string]*Function
}

func (c *Class) Type() ObjType {
	return ClassObj
}

func (c *Class) Inspect() string {
	return "class"
}
