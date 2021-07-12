package interpreter

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/token"
	"fmt"
	"strings"
)

type Interpreter struct {
	program []ast.Stmt
}

func NewInterpreter(program []ast.Stmt) *Interpreter {
	i := &Interpreter{program: program}
	return i
}

func (i *Interpreter) Interpret() interface{} {
	var result interface{}
	for _, stmt := range i.program {
		result = i.evalStmt(stmt)
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
	return expr.Value.Lexeme
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	right := i.evalExpr(expr.Right)
	rType := typeOf(right)
	switch expr.Operator.Type {
	case token.MINUS:
		if rType == 'f' {
			return -right.(float32)
		}
		return "function argument value, type, or count is invalid"
	case token.BANG:
		if rType == 'b' {
			return !right.(bool)
		}
		return "function argument value, type, or count is invalid"
	default:
		return "command contains unrecognized phrase/keyword"
	}
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	left := i.evalExpr(expr.Left)
	right := i.evalExpr(expr.Right)

	lType := typeOf(left)
	rType := typeOf(right)

	ope := expr.Operator.Type

	if lType == 'f' && rType == 'f' {
		return binaryNumber(left.(float32), ope, right.(float32))
	} else if lType == 's' && rType == 's' {
		return binaryString(left.(string), ope, right.(string))
	}
	return fmt.Errorf("operator/operand type mismatch")
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) interface{} {
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	var result interface{}
	for _, s := range stmt.Statements {
		result = i.evalStmt(s)
	}
	return result
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.FunctionStmt) interface{} {
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) interface{} {
	return i.evalExpr(stmt.Value)
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	condition := i.evalExpr(stmt.Condition)
	if typeOf(condition) != 'b' {
		return "data type mismatch"
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
			return "error: division by zero."
		}
		return left / right
	default:
		return "data type mismatch"
	}
}

func binaryString(left string, ope token.TokenType, right string) interface{} {
	switch ope {
	case token.PLUS:
		return left + right
	case token.MINUS:
		return strings.TrimRight(left, " ") + right
	default:
		return "data type mismatch"
	}
}

func typeOf(t interface{}) byte {
	switch t.(type) {
	case string:
		return 's'
	case float32:
		return 'f'
	case bool:
		return 'b'
	default:
		return 'x'
	}
}
