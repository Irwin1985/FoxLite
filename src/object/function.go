package object

import "FoxLite/src/ast"

type Function struct {
	Name       string
	Parameters []*ast.Literal
	Body       *ast.BlockStmt
}

func (f *Function) Type() ObjType {
	return FuncObj
}

func (f *Function) Inspect() string {
	return "func ok"
}
