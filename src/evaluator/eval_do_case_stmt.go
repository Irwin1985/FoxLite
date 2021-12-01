package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"fmt"
)

func evalDoCaseStmt(node *ast.DoCaseStmt, env *object.Environment) object.Object {
	// Comenzamos a recorrer los Branches y nos detenemos si uno es True
	for _, branch := range node.Statements {
		cond := Eval(branch.Condition, env)
		if isError(cond) {
			return cond
		}
		if cond.Type() != object.BooleanObj {
			return object.NewError(fmt.Sprintf("non-bool type `%v` used as case condition", cond.Type()))
		}
		if cond.(*object.Boolean).Value {
			return Eval(branch.Body, env)
		}
	}
	// Revisamos si tenemos una Alternativa
	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return nil
}
