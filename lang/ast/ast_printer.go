package ast

import (
	"FoxLite/lang/token"
	"bytes"
	"fmt"
	"strings"
)

type AstPrinter struct {
	program []Stmt
}

func NewAstPrinter(program []Stmt) *AstPrinter {
	a := &AstPrinter{program: program}
	return a
}

func (a *AstPrinter) PrettyPrint() string {
	var out bytes.Buffer

	for _, stmt := range a.program {
		out.WriteString(fmt.Sprintf("%v\n", a.evalStmt(stmt)))
	}

	return out.String()
}

func (a *AstPrinter) evalStmt(stmt Stmt) interface{} {
	return stmt.Accept(a)
}

func (a *AstPrinter) evalExpr(expr Expr) interface{} {
	return expr.Accept(a)
}

func (a *AstPrinter) VisitExprStmt(stmt *ExprStmt) interface{} {
	return a.evalExpr(stmt.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) interface{} {
	return expr.Value.Lexeme
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return fmt.Sprintf("(%v %v)", expr.Operator.Lexeme, a.evalExpr(expr.Right))
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return fmt.Sprintf("(%v %v %v)", a.evalExpr(expr.Left), expr.Operator.Lexeme, a.evalExpr(expr.Right))
}

func (a *AstPrinter) VisitVarStmt(stmt *VarStmt) interface{} {
	return fmt.Sprintf("%v %v = %v", stmt.Token.Lexeme, a.evalExpr(stmt.Name), a.evalExpr(stmt.Value))
}

func (a *AstPrinter) VisitInlineVarStmt(stmt *InlineVarStmt) interface{} {
	var out bytes.Buffer

	out.WriteString(token.TokenNames[stmt.Scope])
	out.WriteString(" ")
	varList := []string{}
	for _, value := range stmt.Variables {
		v := value.(*VarStmt)
		varList = append(varList, fmt.Sprintf("%v = %v", v.Name.Value.Lexeme.(string), a.evalExpr(v.Value)))
	}
	out.WriteString(strings.Join(varList, ", "))

	return out.String()
}

func (a *AstPrinter) VisitBlockStmt(stmt *BlockStmt) interface{} {
	var out bytes.Buffer
	out.WriteString("\n")
	for _, s := range stmt.Statements {
		out.WriteString(fmt.Sprintf("%v\n", a.evalStmt(s)))
	}
	return out.String()
}

func (a *AstPrinter) VisitFunctionStmt(stmt *FunctionStmt) interface{} {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("FUNCTION %v", stmt.Name.Lexeme))

	if len(stmt.Parameters) > 0 {
		params := []string{}
		for _, param := range stmt.Parameters {
			params = append(params, param.Value.Lexeme.(string))
		}
		out.WriteString(strings.Join(params, ","))
	}
	out.WriteString(fmt.Sprintf("%v", a.evalStmt(stmt.Body)))
	out.WriteString("ENDFUNC")
	return out.String()
}

func (a *AstPrinter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	return fmt.Sprintf("RETURN %v", a.evalExpr(stmt.Value))
}

func (a *AstPrinter) VisitIfStmt(stmt *IfStmt) interface{} {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("IF (%v) THEN", a.evalExpr(stmt.Condition)))
	out.WriteString(fmt.Sprintf("%v", a.evalStmt(stmt.Consequence)))

	if stmt.Alternative != nil {
		out.WriteString(fmt.Sprintf("%v", a.evalStmt(stmt.Alternative)))
	}
	out.WriteString("ENDIF")

	return out.String()
}

func (a *AstPrinter) VisitCallExpr(expr *CallExpr) interface{} {
	return nil
}

func (a *AstPrinter) VisitDoCaseStmt(stmt *DoCaseStmt) interface{} {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("DO CASE\n"))
	for _, b := range stmt.Branches {
		out.WriteString(fmt.Sprintf("   CASE %v\n%v", a.evalExpr(b.Condition), a.evalStmt(b.Consequence)))
	}
	if stmt.Otherwise != nil {
		out.WriteString(fmt.Sprintf("OTHERWISE\n%v", a.evalStmt(stmt.Otherwise)))
	}
	out.WriteString(fmt.Sprintf("ENDCASE"))
	return out.String()
}

func (a *AstPrinter) VisitIifExpr(expr *IifExpr) interface{} {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("IIF(%v, %v, %v)", a.evalExpr(expr.Condition), a.evalExpr(expr.Consequence), a.evalExpr(expr.Alternative)))
	return out.String()
}

func (a *AstPrinter) VisitWhileStmt(stmt *WhileStmt) interface{} {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("DO WHILE (%v)\n", a.evalExpr(stmt.Condition)))
	out.WriteString(fmt.Sprintf("%v", a.evalStmt(stmt.Block)))
	out.WriteString("ENDDO")
	return out.String()
}
