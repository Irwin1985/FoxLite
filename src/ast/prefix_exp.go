package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type PrefixExp struct {
	Op    token.TokenType
	Right Expression
}

func (p *PrefixExp) expressionNode() {}
func (p *PrefixExp) String() string {
	return fmt.Sprintf("%s %v", p.Op, p.Right.String())
}
