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
		if node.Token.Type == token.String {
			return &object.String{Value: val}
		}
		return evalIdentifier(node, env)
	case bool:
		if val {
			return True
		}
		return False
	}
	return Null
}

func evalIdentifier(node *ast.Literal, env *object.Environment) object.Object {
	name := node.Value.(string)
	// resolver el nombre
	result := env.Get(name, true)
	if result == nil {
		lincol := fmt.Sprintf("%d:%d", node.Token.Line, node.Token.Col)
		return &object.Error{Message: fmt.Sprintf("[%s] undefined ident: `%s`", lincol, name)}
	}
	return result
}
