package object

import (
	"FoxLite/lang/token"
	"fmt"
)

type Scope struct {
	Type  token.TokenType
	Value interface{}
}

type Environment struct {
	Store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	e := &Environment{
		Store: make(map[string]interface{}),
	}
	return e
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	e := NewEnvironment()
	e.outer = outer
	return e
}

func (e *Environment) Set(name string, value interface{}, t token.TokenType) {
	s := &Scope{
		Type:  t,
		Value: value,
	}
	e.Store[name] = s
}

func (e *Environment) Assign(name string, value interface{}, t token.TokenType) {
	// try find symbol first
	refEnv := e.GetEnv(name)
	if refEnv == nil { // create the symbol in current env
		e.Set(name, value, t)
		return
	}
	// update the symbol
	scope := refEnv.Store[name].(*Scope)
	scope.Value = value
	refEnv.Store[name] = scope
}

func (e *Environment) Get(name string) (interface{}, error) {
	// find in current environment
	if val, ok := e.Store[name]; ok {
		// if found then return because its a local symbol
		res := val.(*Scope)
		return res.Value, nil
	}
	// if not found then try get the outer environment reference
	if e.outer != nil {
		refEnv := e.outer.GetEnv(name)
		if refEnv != nil {
			scope := refEnv.Store[name].(*Scope)
			// 1. retrieve only PUBLIC and PRIVATE variables
			if scope.Type != token.LOCAL {
				return scope.Value, nil
			}
		}
	}
	return nil, fmt.Errorf(fmt.Sprintf("Variable '%v' is not found.", name))
}

// get the reference of the environment which owns the given symbol
func (e *Environment) GetEnv(name string) *Environment {
	if _, ok := e.Store[name]; ok {
		return e
	}
	if e.outer != nil {
		return e.outer.GetEnv(name)
	}
	return nil
}
