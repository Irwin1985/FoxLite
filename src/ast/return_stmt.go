package ast

import "fmt"

type ReturnStmt struct {
	Value Expression
}

func (r *ReturnStmt) statementNode() {}
func (r *ReturnStmt) String() string {
	return fmt.Sprintf("return %v", r.Value.String())
}
