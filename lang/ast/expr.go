package ast

import "FoxLite/lang/token"

type VisitorExpr interface {
	VisitLiteralExpr(expr *LiteralExpr) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitBinaryExpr(expr *Binary) interface{}
}

type Expr interface {
	Accept(v VisitorExpr) interface{}
}

type LiteralExpr struct {
	Value token.Token
}

func (expr *LiteralExpr) Accept(v VisitorExpr) interface{} {
	return v.VisitLiteralExpr(expr)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (expr *Unary) Accept(v VisitorExpr) interface{} {
	return v.VisitUnaryExpr(expr)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (expr *Binary) Accept(v VisitorExpr) interface{} {
	return v.VisitBinaryExpr(expr)
}
