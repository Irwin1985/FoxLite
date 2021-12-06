package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"fmt"
)

func evalIfExp(node *ast.IfStmt, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	if cond.Type() != object.BooleanObj {
		return object.NewError(fmt.Sprintf("non-bool type `%v` used as if condition", cond.Type()))
	}

	if cond.(*object.Boolean).Value {
		return Eval(node.Consequence, env)
	} else {
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}
	}
	return Null
}
