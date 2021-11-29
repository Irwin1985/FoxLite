package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type ReturnStmt struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStmt) statementNode() {}
func (r *ReturnStmt) String() string {
	return fmt.Sprintf("return %v", r.Value.String())
}
