package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
	"math"
	"strings"
)

func evalArithmeticExp(node *ast.InfixExp, env *object.Environment) object.Object {
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
		return evalBinaryInteger(left.(*object.Integer), right.(*object.Integer), node.Op)
	}
	if lType == object.StringObj && rType == object.StringObj {
		return evalBinaryString(left.(*object.String), right.(*object.String), node.Op)
	}

	// Reportar el error correspondiente
	return reportInfixError(lType, rType)
}

func evalBinaryInteger(left *object.Integer, right *object.Integer, op token.TokenType) object.Object {
	switch op {
	case token.Plus:
		return &object.Integer{Value: left.Value + right.Value}
	case token.Minus:
		return &object.Integer{Value: left.Value - right.Value}
	case token.Mul:
		return &object.Integer{Value: left.Value * right.Value}
	case token.Div:
		if right.Value == 0 {
			return object.NewError("division by zero")
		}
		return &object.Integer{Value: left.Value / right.Value}
	case token.Mod:
		return &object.Integer{Value: math.Mod(left.Value, right.Value)}
	case token.Pow:
		return &object.Integer{Value: math.Pow(left.Value, right.Value)}
	}
	return reportUnexpectedError(op)
}

func evalBinaryString(left *object.String, right *object.String, op token.TokenType) object.Object {
	switch op {
	case token.Plus:
		return &object.String{Value: left.Value + right.Value}
	case token.Minus:
		return &object.String{Value: strings.TrimSpace(left.Value) + right.Value}
	}
	return reportUnexpectedError(op)
}
