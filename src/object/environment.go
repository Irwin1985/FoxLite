package object

type Environment struct {
	storage map[string]Object
	outer   *Environment
}

func NewEnv() *Environment {
	e := &Environment{
		storage: map[string]Object{},
	}
	return e
}

func NewEnclosedEnv(outer *Environment) *Environment {
	e := NewEnv()
	e.outer = outer
	return e
}

func (e *Environment) Set(name string, value Object) Object {
	e.storage[name] = value
	return value
}

func (e *Environment) Get(name string) Object {
	if obj, ok := e.storage[name]; ok {
		return obj
	}
	if e.outer != nil {
		return e.Get(name)
	}
	return nil
}
