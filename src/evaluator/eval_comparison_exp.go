package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
	"fmt"
)

func evalComparisonExp(node *ast.InfixExp, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	lType := left.Type()
	rType := right.Type()

	if lType == object.IntegerObj && rType == object.IntegerObj {
		return evalIntegerComparison(left.(*object.Integer), right.(*object.Integer), node.Op)
	}
	if lType == object.BooleanObj && rType == object.BooleanObj {
		return evalBooleanComparison(left.(*object.Boolean), right.(*object.Boolean), node.Op)
	}
	if lType == object.StringObj && rType == object.StringObj {
		return evalStringComparison(left.(*object.String), right.(*object.String), node.Op)
	}
	return reportInfixError(lType, rType)
}

func evalIntegerComparison(left *object.Integer, right *object.Integer, op token.TokenType) object.Object {
	switch op {
	case token.Less:
		if left.Value < right.Value {
			return True
		}
		return False
	case token.LessEq:
		if left.Value <= right.Value {
			return True
		}
		return False
	case token.Greater:
		if left.Value > right.Value {
			return True
		}
		return False
	case token.GreaterEq:
		if left.Value >= right.Value {
			return True
		}
		return False
	case token.Equal:
		if left.Value == right.Value {
			return True
		}
		return False
	case token.NotEq:
		if left.Value != right.Value {
			return True
		}
		return False
	}
	return reportUnexpectedError(op)
}

func evalBooleanComparison(left *object.Boolean, right *object.Boolean, op token.TokenType) object.Object {
	switch op {
	case token.Equal:
		if left.Value == right.Value {
			return True
		}
		return False
	case token.NotEq:
		if left.Value != right.Value {
			return True
		}
		return False
	}
	return reportUnexpectedError(op)
}

func evalStringComparison(left *object.String, right *object.String, op token.TokenType) object.Object {
	switch op {
	case token.Equal:
		if left.Value == right.Value {
			return True
		}
		return False
	case token.NotEq:
		if left.Value != right.Value {
			return True
		}
		return False
	}
	return object.NewError(fmt.Sprintf("`%s` operator does not support string types", token.GetTokenStr(op)))
}
