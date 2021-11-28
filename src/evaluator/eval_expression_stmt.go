package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalExpressionStmt(node *ast.ExpressionStmt, env *object.Environment) object.Object {
	result := Eval(node.Expression, env)
	if isError(result) {
		return result
	}

	return result
}
