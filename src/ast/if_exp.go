package ast

import (
	"FoxLite/src/token"
	"bytes"
	"fmt"
)

type IfStmt struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (i *IfStmt) statementNode() {}
func (i *IfStmt) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("if %s %s", i.Condition.String(), i.Consequence.String()))
	if i.Alternative != nil {
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}
