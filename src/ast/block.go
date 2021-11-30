package ast

import (
	"bytes"
)

type BlockStmt struct {
	Statements []Statement
}

func (b *BlockStmt) statementNode() {}
func (b *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("\n")
	for _, stmt := range b.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString("\n")
	return out.String()
}
