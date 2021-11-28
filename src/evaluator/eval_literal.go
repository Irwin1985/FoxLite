package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
	"fmt"
)

func evalLiteral(node *ast.Literal, env *object.Environment) object.Object {
	switch val := node.Value.(type) {
	case float64:
		return &object.Integer{Value: val}
	case string:
		if node.Token == token.String {
			return &object.String{Value: val}
		}
		return evalIdentifier(val, env)
	case bool:
		if val {
			return True
		}
		return False
	}
	return Null
}

func evalIdentifier(node string, env *object.Environment) object.Object {
	// resolver el nombre
	result := env.Get(node)
	if result == nil {
		return &object.Error{Message: fmt.Sprintf("undefined ident: `%s`", node)}
	}
	return result
}
