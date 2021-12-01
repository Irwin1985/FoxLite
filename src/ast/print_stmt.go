package ast

import (
	"FoxLite/src/token"
	"bytes"
	"fmt"
)

type PrintStmt struct {
	Token    token.Token
	Messages []Expression
}

func (p *PrintStmt) statementNode() {}
func (p *PrintStmt) String() string {
	var out bytes.Buffer
	for _, msg := range p.Messages {
		out.WriteString(fmt.Sprintf("%s", msg.String()))
	}
	return out.String()
}
