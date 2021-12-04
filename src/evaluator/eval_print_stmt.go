package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"fmt"
)

func evalPrintStmt(node *ast.PrintStmt, env *object.Environment) object.Object {
	// Mostramos por consola todas las expresiones
	for _, exp := range node.Messages {
		res := Eval(exp, env)
		if isError(res) {
			return res
		}
		fmt.Printf("%s ", res.Inspect())
	}
	fmt.Print("\n")

	return None
}
