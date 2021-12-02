package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalReturnStmt(node *ast.ReturnStmt, env *object.Environment) object.Object {
	var val object.Object
	if node.Value != nil {
		val = Eval(node.Value, env)
		if isError(val) {
			return val
		}
	} else {
		val = &object.None{}
	}
	return &object.Return{Value: val}
}
