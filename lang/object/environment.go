package object

import "fmt"

type Environment struct {
	store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	e := &Environment{
		store: make(map[string]interface{}),
	}
	return e
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	e := NewEnvironment()
	e.outer = outer
	return e
}

func (e *Environment) Set(name string, value interface{}) interface{} {
	e.store[name] = value
	return value
}

func (e *Environment) Get(name string) (interface{}, error) {
	if v, ok := e.store[name]; ok {
		return v, nil
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, fmt.Errorf("variable not found [%v]", name)
}
