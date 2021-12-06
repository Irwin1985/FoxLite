package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
)

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	name := node.Name.String()
	f := &object.Function{
		Name:       node.Name.String(),
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
	return env.Set(name, 'g', f)
}
