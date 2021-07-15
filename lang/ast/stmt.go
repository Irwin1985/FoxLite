package ast

import "FoxLite/lang/token"

type VisitorStmt interface {
	VisitBlockStmt(stmt *BlockStmt) interface{}
	VisitExprStmt(stmt *ExprStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
	VisitFunctionStmt(stmt *FunctionStmt) interface{}
	VisitReturnStmt(stmt *ReturnStmt) interface{}
	VisitIfStmt(stmt *IfStmt) interface{}
	VisitInlineVarStmt(stmt *InlineVarStmt) interface{}
	VisitDoCaseStmt(stmt *DoCaseStmt) interface{}
	VisitWhileStmt(stmt *WhileStmt) interface{}
}

type Stmt interface {
	Accept(v VisitorStmt) interface{}
}

type BlockStmt struct {
	Statements []Stmt
}

func (stmt *BlockStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitBlockStmt(stmt)
}

type FunctionStmt struct {
	Name       token.Token
	Parameters []LiteralExpr
	Body       *BlockStmt
}

func (stmt *FunctionStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitFunctionStmt(stmt)
}

type ReturnStmt struct {
	Value Expr
}

func (stmt *ReturnStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitReturnStmt(stmt)
}

type ExprStmt struct {
	Expression Expr
}

func (stmt *ExprStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitExprStmt(stmt)
}

type VarStmt struct {
	Token token.Token
	Name  *LiteralExpr
	Type  token.TokenType
	Value Expr
}

func (stmt *VarStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitVarStmt(stmt)
}

type InlineVarStmt struct {
	Scope     token.TokenType
	Variables []Stmt
}

func (stmt *InlineVarStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitInlineVarStmt(stmt)
}

type IfStmt struct {
	Condition   Expr
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (stmt *IfStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitIfStmt(stmt)
}

type DoCaseStmt struct {
	Branches  []*IfStmt
	Otherwise *BlockStmt
}

func (stmt *DoCaseStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitDoCaseStmt(stmt)
}

type WhileStmt struct {
	Condition Expr
	Block     *BlockStmt
}

func (stmt *WhileStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitWhileStmt(stmt)
}
