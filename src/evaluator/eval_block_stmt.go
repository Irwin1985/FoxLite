package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalBlockStmt(node *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
		rType := result.Type()
		if rType == object.ReturnObj || rType == object.ErrorObj || rType == object.ExitObj || rType == object.LoopObj {
			return result
		}
	}
	return result
}
