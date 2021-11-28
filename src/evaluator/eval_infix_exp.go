package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
)

// evalInfixExp => evalúa las expresiones infijas que pueden ser:
// +, -, *, /, %, ^, ==, !=, <, <=, >, >=, and, or
func evalInfixExp(node *ast.InfixExp, env *object.Environment) object.Object {
	switch node.Op {
	case token.And, token.Or:
		return evalLogicalExp(node, env)
	case token.Plus, token.Minus, token.Mul, token.Div, token.Mod, token.Pow:
		return evalArithmeticExp(node, env)
	case token.Less, token.LessEq, token.Greater, token.GreaterEq, token.Equal, token.NotEq:
		return evalComparisonExp(node, env)
	}
	return reportUnexpectedError(node.Op)
}
