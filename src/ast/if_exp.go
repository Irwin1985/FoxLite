package ast

import (
	"FoxLite/src/token"
	"bytes"
	"fmt"
)

type IfExp struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (i *IfExp) expressionNode() {}
func (i *IfExp) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("if %s %s", i.Condition.String(), i.Consequence.String()))
	if i.Alternative != nil {
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}
