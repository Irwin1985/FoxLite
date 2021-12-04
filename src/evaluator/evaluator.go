package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"FoxLite/src/token"
	"fmt"
)

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}
var None = &object.None{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStmt:
		return evalBlockStmt(node, env)
	case *ast.ExpressionStmt:
		return evalExpressionStmt(node, env)
	case *ast.Literal:
		return evalLiteral(node, env)
	case *ast.ReturnStmt:
		return evalReturnStmt(node, env)
	case *ast.InfixExp:
		return evalInfixExp(node, env)
	case *ast.VarStmt:
		return evalVarStmt(node, env)
	case *ast.IfExp:
		return evalIfExp(node, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.PrintStmt:
		return evalPrintStmt(node, env)
	case *ast.DoCaseStmt:
		return evalDoCaseStmt(node, env)
	case *ast.While:
		return evalWhileStmt(node, env)
	case *ast.Loop:
		return &object.Loop{}
	case *ast.Exit:
		return &object.Exit{}
	case *ast.Input:
		return evalInputStmt(node, env)
	default:
		return None
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}
	return false
}

func reportInfixError(lType object.ObjType, rType object.ObjType) object.Object {
	if lType == object.StringObj || lType == object.IntegerObj {
		return object.NewError(fmt.Sprintf("infix expr: cannot use `%s` (right expression) as `%s`", object.TypeToStr(rType), object.TypeToStr(lType)))
	}
	if lType == object.BooleanObj {
		return object.NewError("bool types only have the following operators defined: `!`, `==`, `!=`, `or`, `and`")
	}
	return Null
}

func reportUnexpectedError(op token.TokenType) object.Object {
	return object.NewError(fmt.Sprintf("unexpected token `%s`", token.GetTokenStr(op)))
}
