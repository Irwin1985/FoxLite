package ast

import "FoxLite/src/token"

type CaseBranch struct {
	Condition Expression
	Body      *BlockStmt
}

type DoCaseStmt struct {
	Token       token.Token
	Statements  []*CaseBranch
	Alternative *BlockStmt
}

func (d *DoCaseStmt) statementNode() {}
func (d *DoCaseStmt) String() string {
	return "do case"
}
