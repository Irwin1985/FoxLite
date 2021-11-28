package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
)

func evalPrefixExp(node *ast.PrefixExp, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	switch node.Op {
	case token.Not:
		if right.Type() != object.BooleanObj {
			return object.NewError("! operator can only be used with bool types")
		}
		val := right.(*object.Boolean).Value
		if val == true {
			return False
		}
		return True
	case token.Minus:
		if right.Type() != object.IntegerObj {
			return object.NewError("- operator can only be used with numeric types")
		}
		return &object.Integer{Value: right.(*object.Integer).Value * -1}
	}
	return reportUnexpectedError(node.Op)
}
