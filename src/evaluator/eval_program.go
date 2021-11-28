package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Eval(stmt, env)
		// En caso de ser Return devolvemos su valor
		// ya que estamos a nivel de un programa (prg -> flp)
		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}
