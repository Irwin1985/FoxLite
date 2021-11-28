package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalReturnStmt(node *ast.ReturnStmt, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	return &object.Return{Value: val}
}
