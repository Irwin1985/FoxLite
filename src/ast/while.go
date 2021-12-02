package ast

import "FoxLite/src/token"

type While struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStmt
}

func (w *While) statementNode() {}
func (w *While) String() string {
	return "while"
}
