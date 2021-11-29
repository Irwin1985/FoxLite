package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalVarStmt(node *ast.VarStmt, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	return env.Set(node.Name, node.Scope, val)
}
