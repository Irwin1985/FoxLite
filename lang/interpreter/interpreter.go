package interpreter

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/object"
	"FoxLite/lang/token"
	"fmt"
	"strings"
)

type Interpreter struct {
	env       *object.Environment
	globalEnv *object.Environment
	program   []ast.Stmt
}

func NewInterpreter(program []ast.Stmt, env *object.Environment) *Interpreter {
	i := &Interpreter{
		env:       env,
		globalEnv: env,
		program:   program,
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
			return result.(*object.Return).Value
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
		v, err := i.env.Get(expr.Value.Lexeme.(string))
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

	if ope == token.AND || ope == token.OR {
		return i.evalLogicalExpr(expr)
	}

	if ope == token.PLUS_EQ || ope == token.MINUS_EQ || ope == token.MUL_EQ || ope == token.DIV_EQ {
		return i.evalShortAssignment(expr)
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
	value := i.evalExpr(stmt.Value)
	if isError(value) {
		return value
	}
	// PUBLIC variables goes directly in globalEnv
	if stmt.Token.Type == token.PUBLIC {
		i.globalEnv.Set(name, value, stmt.Token.Type)
		return nil
	}

	i.env.Set(name, value, stmt.Token.Type)
	return nil
}

func (i *Interpreter) VisitInlineVarStmt(stmt *ast.InlineVarStmt) interface{} {
	curEnv := i.env
	// PUBLIC variables goes directly in globalEnv
	if stmt.Scope == token.PUBLIC {
		curEnv = i.globalEnv
	}
	for _, value := range stmt.Variables {
		node := value.(*ast.VarStmt)
		value := i.evalExpr(node.Value)
		if isError(value) {
			return value
		}
		curEnv.Set(node.Name.Value.Lexeme.(string), value, stmt.Scope)
	}
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
	funObj.Env = i.env
	name := funObj.Name
	_, err := i.env.Get(name)
	if err == nil {
		return newError(fmt.Sprintf("invalid function name: '%s'", name))
	}
	i.env.Set(name, funObj, token.PRIVATE)

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
		args := i.evalArguments(expr.Arguments)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return i.applyFunction(funObj, args)
	}
	return nil
}

func (i *Interpreter) evalArguments(args []ast.Expr) []interface{} {
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
	if lit, ok := expr.Left.(*ast.LiteralExpr); ok {
		if name, ok := lit.Value.Lexeme.(string); ok {
			result := i.evalExpr(expr.Right)
			i.env.Assign(name, result, token.PRIVATE)
			return nil
		}
	}
	return newError("syntax error: left hand side must be an identifier")
}

func (i *Interpreter) evalLogicalExpr(expr *ast.Binary) interface{} {
	left := i.evalExpr(expr.Left)
	if isError(left) {
		return left
	}
	if typeOf(left) != 'b' {
		return newError("function argument value, type or count is invalid.")
	}
	leftVal := left.(bool)
	ope := expr.Operator.Type

	if ope == token.AND && !leftVal {
		return i.evalRightExpr(expr.Right)
	}
	if ope == token.OR && !leftVal {
		return i.evalRightExpr(expr.Right)
	}
	return true
}

func (i *Interpreter) evalRightExpr(expr ast.Expr) interface{} {
	right := i.evalExpr(expr)
	if isError(right) {
		return right
	}
	if typeOf(right) != 'b' {
		return newError("function argument value, type or count is invalid.")
	}
	return right.(bool)
}

func (i *Interpreter) evalShortAssignment(expr *ast.Binary) interface{} {
	if left, ok := expr.Left.(*ast.LiteralExpr); ok {
		name := left.Value.Lexeme.(string)
		refEnv := i.env.GetEnv(name)
		if refEnv == nil {
			return newError(fmt.Sprintf("Variable '%v' is not found.", name))
		}
		scope := refEnv.Store[name].(*object.Scope)
		leftValue := scope.Value

		if typeOf(leftValue) == 'n' || typeOf(leftValue) == 's' {
			right := i.evalExpr(expr.Right)
			if isError(right) {
				return right
			}
			ope := expr.Operator.Type
			switch right.(type) {
			case string:
				if typeOf(right) != 's' {
					return newError("Operator/operand type mismatch.")
				}
				switch ope {
				case token.PLUS_EQ:
					res := leftValue.(string) + right.(string)
					scope.Value = res
					refEnv.Store[name] = scope
				case token.MINUS_EQ:
					res := strings.TrimRight(leftValue.(string), " ") + right.(string)
					scope.Value = res
					refEnv.Store[name] = scope
				default:
					return newError("function argument value, type or count is invalid.")
				}
			case float32:
				if typeOf(right) != 'n' {
					return newError("Operator/operand type mismatch.")
				}
				switch ope {
				case token.PLUS_EQ:
					res := leftValue.(float32) + right.(float32)
					scope.Value = res
					refEnv.Store[name] = scope
				case token.MINUS_EQ:
					res := leftValue.(float32) - right.(float32)
					scope.Value = res
					refEnv.Store[name] = scope
				case token.MUL_EQ:
					res := leftValue.(float32) * right.(float32)
					scope.Value = res
					refEnv.Store[name] = scope
				case token.DIV_EQ:
					rightVal := scope.Value.(float32)
					if rightVal == 0 {
						return newError("division by zero.")
					}
					res := leftValue.(float32) / right.(float32)
					scope.Value = res
					refEnv.Store[name] = scope
				default:
					return newError("function argument value, type or count is invalid.")
				}
			}
			return nil
		}
	}
	return newError("function argument value, type or count is invalid.")
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
	case token.LT:
		return left < right
	case token.GT:
		return left > right
	case token.LEQ:
		return left <= right
	case token.GEQ:
		return left >= right
	case token.EQ:
		return left == right
	case token.NEQ:
		return left != right
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
	// do some parameters validation
	fnParNum := len(fn.Parameters)
	argNum := len(args)
	// 1. function does not implement parameters and got some.
	if fnParNum == 0 && argNum > 0 {
		return newError("No PARAMETER statement is found")
	}
	// 2. function expect more parameters that expected ones.
	if argNum > fnParNum {
		return newError("Must specify additional parameters.")
	}
	extendedEnv := extendFunctionEnv(fn, args)

	oldEnv := i.env
	i.env = extendedEnv

	evaluated := i.evalStmt(fn.Body)
	i.env = oldEnv
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []interface{}) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value.Lexeme.(string), args[paramIdx], token.PRIVATE)
		} else {
			// add false as default
			env.Set(param.Value.Lexeme.(string), false, token.PRIVATE)
		}
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
	case *object.Error:
		return true
	default:
		return false
	}
}
