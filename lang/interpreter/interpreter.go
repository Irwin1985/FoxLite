package interpreter

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/object"
	"FoxLite/lang/token"
	"fmt"
	"strings"
)

type Interpreter struct {
	Env     *object.Environment
	program []ast.Stmt
}

func NewInterpreter(program []ast.Stmt) *Interpreter {
	i := &Interpreter{
		Env:     object.NewEnvironment(),
		program: program,
	}
	return i
}

func (i *Interpreter) Interpret() interface{} {
	var result interface{}
	for _, stmt := range i.program {
		result = i.evalStmt(stmt)
		if isError(result) {
			return result
		} else if typeOf(result) == 'r' {
			return result.(object.Return).Value
		}
	}
	return result
}

func (i *Interpreter) evalStmt(stmt ast.Stmt) interface{} {
	return stmt.Accept(i)
}

func (i *Interpreter) evalExpr(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	return i.evalExpr(stmt.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) interface{} {
	if expr.Value.Type == token.IDENT {
		v, err := i.Env.Get(expr.Value.Lexeme.(string))
		if err != nil {
			return newError(fmt.Sprintf("%v", err))
		}
		return v
	}
	return expr.Value.Lexeme
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	right := i.evalExpr(expr.Right)
	if isError(right) {
		return right
	}
	rType := typeOf(right)
	switch expr.Operator.Type {
	case token.MINUS:
		if rType == 'n' {
			return -right.(float32)
		}
		return newError("function argument value, type, or count is invalid")
	case token.BANG:
		if rType == 'b' {
			return !right.(bool)
		}
		return newError("function argument value, type, or count is invalid")
	default:
		return newError("command contains unrecognized phrase/keyword")
	}
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	ope := expr.Operator.Type

	if ope == token.ASSIGN {
		return i.evalAssignment(expr)
	}

	left := i.evalExpr(expr.Left)
	if isError(left) {
		return left
	}
	right := i.evalExpr(expr.Right)
	if isError(right) {
		return right
	}
	lType := typeOf(left)
	rType := typeOf(right)

	if lType == 'n' && rType == 'n' {
		return binaryNumber(left.(float32), ope, right.(float32))
	} else if lType == 's' && rType == 's' {
		return binaryString(left.(string), ope, right.(string))
	}
	return newError("operator/operand type mismatch")
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) interface{} {
	name := stmt.Name.Value.Lexeme.(string)
	i.Env.Set(name, i.evalExpr(stmt.Value))
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	var result interface{}
	for _, s := range stmt.Statements {
		result = i.evalStmt(s)
		if isError(result) || typeOf(result) == 'r' {
			return result
		}
	}
	return result
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) interface{} {
	funObj := &object.Function{}
	funObj.Name = stmt.Name.Lexeme.(string)
	funObj.Parameters = stmt.Parameters
	funObj.Body = stmt.Body
	funObj.Env = i.Env

	i.Env.Set(funObj.Name, funObj)

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) interface{} {
	left := i.evalExpr(expr.Function)
	if isError(left) {
		return left
	}
	if typeOf(left) != 'f' {
		return newError("Syntax error: left hand side must be a function definition")
	}
	if funObj, ok := left.(*object.Function); ok {
		args := i.evalExpressions(expr.Arguments)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return i.applyFunction(funObj, args)
	}
	return nil
}

func (i *Interpreter) evalExpressions(args []ast.Expr) []interface{} {
	exps := []interface{}{}

	var result interface{}
	for _, e := range args {
		result = i.evalExpr(e)
		if result == nil {
			return []interface{}{e}
		}
		exps = append(exps, result)
	}

	return exps
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) interface{} {
	result := i.evalExpr(stmt.Value)
	if isError(result) {
		return result
	}

	return &object.Return{Value: result}
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	condition := i.evalExpr(stmt.Condition)
	if isError(condition) {
		return condition
	}
	if typeOf(condition) != 'b' {
		return newError("data type mismatch")
	}
	if condition.(bool) {
		return i.evalStmt(stmt.Consequence)
	} else {
		if stmt.Alternative != nil {
			return i.evalStmt(stmt.Alternative)
		}
	}
	return nil
}

// HELPER FUNCTIONS
func (i *Interpreter) evalAssignment(expr *ast.Binary) interface{} {
	leftVal := i.evalExpr(expr.Left)
	if isError(leftVal) {
		return leftVal
	}
	if name, ok := leftVal.(string); ok {
		return i.Env.Set(name, i.evalExpr(expr.Right))
	}
	return newError("syntax error: left hand side must be an identifier")
}

func binaryNumber(left float32, ope token.TokenType, right float32) interface{} {
	switch ope {
	case token.PLUS:
		return left + right
	case token.MINUS:
		return left - right
	case token.MUL:
		return left * right
	case token.DIV:
		if right == 0 {
			return newError("error: division by zero")
		}
		return left / right
	default:
		return newError("data type mismatch")
	}
}

func binaryString(left string, ope token.TokenType, right string) interface{} {
	switch ope {
	case token.PLUS:
		return left + right
	case token.MINUS:
		return strings.TrimRight(left, " ") + right
	default:
		return newError("data type mismatch")
	}
}

func (i *Interpreter) applyFunction(fn *object.Function, args []interface{}) interface{} {
	extendedEnv := extendFunctionEnv(fn, args)

	oldEnv := i.Env
	i.Env = extendedEnv

	evaluated := i.evalStmt(fn.Body)
	i.Env = oldEnv
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []interface{}) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value.Lexeme.(string), args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj interface{}) interface{} {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}
	return obj
}

func typeOf(t interface{}) byte {
	switch t.(type) {
	case string:
		return 's'
	case float32:
		return 'n'
	case bool:
		return 'b'
	case *object.Function:
		return 'f'
	case *object.Return:
		return 'r'
	default:
		return 'x'
	}
}

func newError(msg string) *object.Error {
	return &object.Error{Message: msg}
}

func isError(i interface{}) bool {
	if i == nil {
		return false
	}

	switch i.(type) {
	case object.Error:
		return true
	default:
		return false
	}
}
