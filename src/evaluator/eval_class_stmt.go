package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalClassStmt(node *ast.Class, env *object.Environment) object.Object {
	class := &object.Class{
		Name:       node.Name,
		Properties: map[string]object.Object{},
		Methods:    map[string]*object.Function{},
	}

	// Evaluate all properties
	for key, prop := range node.Properties {
		res := Eval(prop, env)
		if isError(res) {
			return res
		}
		class.Properties[key] = res
	}

	// Evaluate all methods
	for key, fn := range node.Methods {
		res := Eval(fn, env)
		if isError(res) {
			return res
		}
		class.Methods[key] = res.(*object.Function)
	}

	return class
}
