package object

import "FoxLite/lang/ast"

type Function struct {
	Name       string
	Env        *Environment
	Parameters []ast.LiteralExpr
	Body       *ast.BlockStmt
}

type Error struct {
	Message string
}

type Return struct {
	Value interface{}
}
