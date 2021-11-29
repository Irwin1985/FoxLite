package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type ExpressionStmt struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStmt) statementNode() {}
func (e *ExpressionStmt) String() string {
	return fmt.Sprintf("%s\n", e.Expression.String())
}
