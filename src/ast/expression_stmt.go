package ast

import "fmt"

type ExpressionStmt struct {
	Expression Expression
}

func (e *ExpressionStmt) statementNode() {}
func (e *ExpressionStmt) String() string {
	return fmt.Sprintf("%s\n", e.Expression.String())
}
