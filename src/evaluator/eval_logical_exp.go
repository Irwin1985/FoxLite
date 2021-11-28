package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
	"fmt"
)

func evalLogicalExp(node *ast.InfixExp, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	op := "and"
	if node.Op == token.Or {
		op = "or"
	}
	if left.Type() != object.BooleanObj {
		return object.NewError(fmt.Sprintf("left operand for `%s` is not a boolean", op))
	}
	switch node.Op {
	case token.And:
		if !left.(*object.Boolean).Value {
			return False
		}
		return evalRightExp(node, env)
	case token.Or:
		if left.(*object.Boolean).Value {
			return True
		}
		return evalRightExp(node, env)
	}
	return Null
}

func evalRightExp(node *ast.InfixExp, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	op := "and"
	if node.Op == token.Or {
		op = "or"
	}
	if right.Type() != object.BooleanObj {
		return object.NewError(fmt.Sprintf("right operand for `%s` is not a boolean", op))
	}
	if right.(*object.Boolean).Value {
		return True
	}
	return False
}
