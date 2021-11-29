package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type VarStmt struct {
	Token token.Token
	Scope byte
	Type  byte
	Name  string
	Value Expression
}

func (v *VarStmt) statementNode() {}

func (v *VarStmt) String() string {
	return fmt.Sprintf("%s = %s", v.Name, v.Value.String())
}
