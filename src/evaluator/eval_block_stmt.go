package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalBlockStmt(node *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == object.ReturnObj || result.Type() == object.ErrorObj {
				return result
			}
		}
	}
	return result
}
