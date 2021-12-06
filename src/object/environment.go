package object

type Vector struct {
	Scope byte
	Value Object
}

type Environment struct {
	storage map[string]*Vector
	outer   *Environment
}

func NewEnv() *Environment {
	e := &Environment{
		storage: map[string]*Vector{},
		outer:   nil,
	}
	return e
}

func NewEnclosedEnv(outer *Environment) *Environment {
	e := NewEnv()
	e.outer = outer
	return e
}

func (e *Environment) Set(name string, scope byte, value Object) Object {
	// Creamos un nuevo vector
	v := &Vector{
		Scope: scope,
		Value: value,
	}
	e.storage[name] = v
	return value
}

func (e *Environment) Get(name string, outCall bool) Object {
	if vec, ok := e.storage[name]; ok {
		if outCall { // si llaman desde afuera: devolvemos sin validar scope
			return vec.Value
		} else {
			// si es una llamada recursiva (interna): devolvemos solo private y public
			if vec.Scope == 'p' || vec.Scope == 'g' {
				return vec.Value
			}
			return nil
		}
	} else {
		if e.outer != nil {
			return e.outer.Get(name, false)
		}
	}
	return nil
}
