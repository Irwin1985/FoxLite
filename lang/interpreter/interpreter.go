package interpreter

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/object"
	"FoxLite/lang/token"
	"fmt"
	"strings"
)

type Interpreter struct {
	env  *object.Environment
	gEnv *object.Environment
	prog []ast.Stmt
}

func NewInterpreter(prog []ast.Stmt, env *object.Environment) *Interpreter {
	i := &Interpreter{
		env:  env,
		gEnv: env,
		prog: prog,
	}
	return i
}

func (i *Interpreter) Interpret() interface{} {
	var r interface{}
	for _, st := range i.prog {
		r = i.evalStmt(st)
		if typeOf(r) == 'e' {
			return r
		}
		if o, ok := r.(*object.Return); ok {
			return o.Value
		}
	}
	return r
}

func (i *Interpreter) evalStmt(stmt ast.Stmt) interface{} {
	return stmt.Accept(i)
}

func (i *Interpreter) evalExpr(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	var r interface{}
	for _, s := range stmt.Statements {
		r = i.evalStmt(s)
		if o, ok := r.(*object.Return); ok {
			return o
		}
		if typeOf(r) == 'e' {
			return r
		}
	}
	return nil
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
	v := i.evalExpr(expr.Right)
	switch expr.Operator.Type {
	case token.MINUS:
		if r, ok := v.(float32); ok {
			return -r
		}
	case token.BANG:
		if r, ok := v.(bool); ok {
			return !r
		}
	default:
		return newError("command contains unrecognized phrase/keyword")
	}
	return newError("function argument value, type, or count is invalid")
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	op := expr.Operator.Type

	if op == token.ASSIGN {
		return i.evalAssignment(expr)
	}

	if op == token.AND || op == token.OR {
		return i.evalLogicalExpr(expr)
	}

	if op == token.PLUS_EQ || op == token.MINUS_EQ || op == token.MUL_EQ || op == token.DIV_EQ {
		return i.evalShortAssignment(expr)
	}

	l := i.evalExpr(expr.Left)
	r := i.evalExpr(expr.Right)
	lt := typeOf(l)
	rt := typeOf(r)

	if lt == rt {
		if lt == 'n' {
			return binaryNumber(l.(float32), op, r.(float32))
		} else if lt == 's' {
			return binaryString(l.(string), op, r.(string))
		}
	}
	return newError("operator/operand type mismatch")
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) interface{} {
	n := stmt.Name.Value.Lexeme.(string)
	var value interface{}
	if stmt.Value != nil {
		value = i.evalExpr(stmt.Value)
		if typeOf(value) == 'e' {
			return value
		}
	} else {
		switch stmt.Type {
		case token.STRING_T:
			value = ""
		case token.NUMBER_T:
			value = 0
		case token.BOOLEAN_T:
			value = false
		}
	}
	// PUBLIC variables goes directly in globalEnv
	if stmt.Token.Type == token.PUBLIC {
		i.gEnv.Set(n, value, stmt.Token.Type)
		return nil
	}

	i.env.Set(n, value, stmt.Token.Type)
	return nil
}

func (i *Interpreter) VisitInlineVarStmt(stmt *ast.InlineVarStmt) interface{} {
	curEnv := i.env
	// PUBLIC variables goes directly in globalEnv
	if stmt.Scope == token.PUBLIC {
		curEnv = i.gEnv
	}
	for _, vars := range stmt.Variables {
		node := vars.(*ast.VarStmt)
		r := i.evalExpr(node.Value)
		if typeOf(r) == 'e' {
			return r
		}
		curEnv.Set(node.Name.Value.Lexeme.(string), r, stmt.Scope)
	}
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) interface{} {
	fn := &object.Function{}
	fn.Name = stmt.Name.Lexeme.(string)
	fn.Parameters = stmt.Parameters
	fn.Body = stmt.Body
	fn.Env = i.env
	name := fn.Name
	_, err := i.env.Get(name)
	if err == nil {
		return newError(fmt.Sprintf("name in use: '%s'", name))
	}
	i.env.Set(name, fn, token.PRIVATE)

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) interface{} {
	r := i.evalExpr(expr.Function)
	if typeOf(r) != 'f' {
		return newError("Syntax error: left hand side must be a function definition")
	}
	fn := r.(*object.Function)
	args := i.evalArguments(expr.Arguments)
	if len(args) == 1 && typeOf(args[0]) == 'e' {
		return args[0]
	}

	return i.applyFunction(fn, args)
}

func (i *Interpreter) evalArguments(args []ast.Expr) []interface{} {
	exps := []interface{}{}

	var r interface{}
	for _, e := range args {
		r = i.evalExpr(e)
		if typeOf(r) == 'e' {
			return []interface{}{r}
		}
		exps = append(exps, r)
	}

	return exps
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) interface{} {
	r := i.evalExpr(stmt.Value)
	if typeOf(r) == 'e' {
		return r
	}

	return &object.Return{Value: r}
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	r := i.evalExpr(stmt.Condition)
	if typeOf(r) != 'b' {
		return newError("data type mismatch")
	}
	if b, ok := r.(bool); ok && b {
		return i.evalStmt(stmt.Consequence)
	} else {
		if stmt.Alternative != nil {
			return i.evalStmt(stmt.Alternative)
		}
	}
	return nil
}

func (i *Interpreter) VisitIifExpr(expr *ast.IifExpr) interface{} {
	r := i.evalExpr(expr.Condition)
	if typeOf(r) != 'b' {
		return newError("data type mismatch")
	}
	if b, ok := r.(bool); ok && b {
		return i.evalExpr(expr.Consequence)
	} else {
		return i.evalExpr(expr.Alternative)
	}
}

func (i *Interpreter) VisitDoCaseStmt(stmt *ast.DoCaseStmt) interface{} {
	for _, bl := range stmt.Branches {
		r := i.evalExpr(bl.Condition)
		if typeOf(r) != 'b' {
			return newError("data type mismatch")
		}
		if b, ok := r.(bool); ok && b {
			return i.evalStmt(bl.Consequence)
		}
	}
	if stmt.Otherwise != nil {
		return i.evalStmt(stmt.Otherwise)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) interface{} {
	for {
		if b, ok := i.evalExpr(stmt.Condition).(bool); ok && b {
			r := i.evalStmt(stmt.Block)
			if o, ok := r.(*object.Return); ok {
				return o
			}
			if typeOf(r) == 'e' {
				return r
			}
		} else {
			return newError("data type mismatch")
		}
	}
}

// HELPER FUNCTIONS
func (i *Interpreter) evalAssignment(expr *ast.Binary) interface{} {
	if l, ok := expr.Left.(*ast.LiteralExpr); ok {
		if name, ok := l.Value.Lexeme.(string); ok {
			r := i.evalExpr(expr.Right)
			i.env.Assign(name, r, token.PRIVATE)
			return nil
		}
	}
	return newError("syntax error: left hand side must be an identifier")
}

func (i *Interpreter) evalLogicalExpr(expr *ast.Binary) interface{} {
	r := i.evalExpr(expr.Left)
	if v, ok := r.(bool); ok {
		op := expr.Operator.Type

		if op == token.AND {
			if v {
				return i.evalRightExpr(expr.Right)
			}
			return false
		}
		if op == token.OR {
			if !v {
				return i.evalRightExpr(expr.Right)
			}
			return true
		}
		return nil
	}
	return newError("function argument value, type or count is invalid.")
}

func (i *Interpreter) evalRightExpr(expr ast.Expr) interface{} {
	if v, ok := i.evalExpr(expr).(bool); ok {
		return v
	}
	return newError("function argument value, type or count is invalid.")
}

func (i *Interpreter) evalShortAssignment(expr *ast.Binary) interface{} {
	if l, ok := expr.Left.(*ast.LiteralExpr); ok {
		n := l.Value.Lexeme.(string)
		refEnv := i.env.GetEnv(n)
		if refEnv == nil {
			return newError(fmt.Sprintf("Variable '%v' is not found.", n))
		}

		scope := refEnv.Store[n].(*object.Scope)
		lv := scope.Value
		t := typeOf(lv)
		if t != 'n' && t != 's' {
			return newError("operator/operand type mismatch.")
		}

		rv := i.evalExpr(expr.Right)
		rt := typeOf(rv)
		op := expr.Operator.Type

		switch rt {
		case 's':
			var lv = lv.(string)
			var rv = rv.(string)
			switch op {
			case token.PLUS_EQ:
				scope.Value = lv + rv
				refEnv.Store[n] = scope
			case token.MINUS_EQ:
				scope.Value = strings.TrimRight(lv, " ") + rv
				refEnv.Store[n] = scope
			default:
				return newError("missing operand")
			}
		case 'n':
			var lv = lv.(float32)
			var rv = rv.(float32)
			switch op {
			case token.PLUS_EQ:
				scope.Value = lv + rv
				refEnv.Store[n] = scope
			case token.MINUS_EQ:
				scope.Value = lv - rv
				refEnv.Store[n] = scope
			case token.MUL_EQ:
				scope.Value = lv * rv
				refEnv.Store[n] = scope
			case token.DIV_EQ:
				if rv == 0 {
					return newError("division by zero.")
				}
				scope.Value = lv / rv
				refEnv.Store[n] = scope
			default:
				return newError("missing operand")
			}
		default:
			return newError("operator/operand type mismatch.")
		}
		return nil
	}
	return newError("function argument value, type or count is invalid.")
}

func binaryNumber(lv float32, op token.TokenType, rv float32) interface{} {
	switch op {
	case token.PLUS:
		return lv + rv
	case token.MINUS:
		return lv - rv
	case token.MUL:
		return lv * rv
	case token.DIV:
		if rv == 0 {
			return newError("division by zero")
		}
		return lv / rv
	case token.LT:
		return lv < rv
	case token.GT:
		return lv > rv
	case token.LEQ:
		return lv <= rv
	case token.GEQ:
		return lv >= rv
	case token.EQ:
		return lv == rv
	case token.NEQ:
		return lv != rv
	default:
		return newError("missing operand")
	}
}

func binaryString(lv string, op token.TokenType, rv string) interface{} {
	switch op {
	case token.PLUS:
		return lv + rv
	case token.MINUS:
		return strings.TrimRight(lv, " ") + rv
	case token.EQ:
		return lv == rv
	default:
		return newError("missing operand")
	}
}

func (i *Interpreter) applyFunction(fn *object.Function, args []interface{}) interface{} {
	// do some parameters validation
	fnPn := len(fn.Parameters)
	argn := len(args)
	// 1. function does not implement parameters and got some.
	if fnPn == 0 && argn > 0 {
		return newError("No PARAMETER statement is found")
	}
	// 2. function expect more parameters that expected ones.
	if argn > fnPn {
		return newError("Must specify additional parameters.")
	}
	newEnv := extendFunctionEnv(fn, args)

	oldEnv := i.env
	i.env = newEnv

	evaluated := i.evalStmt(fn.Body)
	i.env = oldEnv

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []interface{}) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, p := range fn.Parameters {
		if i < len(args) {
			env.Set(p.Value.Lexeme.(string), args[i], token.PRIVATE)
		} else {
			env.Set(p.Value.Lexeme.(string), false, token.PRIVATE) // false as default
		}
	}

	return env
}

func unwrapReturnValue(obj interface{}) interface{} {
	if rv, ok := obj.(*object.Return); ok {
		return rv.Value
	}
	return obj
}

func typeOf(t interface{}) byte {
	if t == nil {
		return 'x'
	}
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
	case *object.Error:
		return 'e'
	default:
		return 'u'
	}
}

func newError(msg string) *object.Error {
	return &object.Error{Message: msg}
}
