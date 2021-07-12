package object

import "FoxLite/lang/ast"

type FunctionObj struct {
	Parameters []ast.LiteralExpr
	Body       *ast.BlockStmt
}
