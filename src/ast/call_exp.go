package ast

import "FoxLite/src/token"

type CallExp struct {
	Token  token.Token
	Caller Expression
	Args   []Expression
}

func (c *CallExp) expressionNode() {}
func (c *CallExp) String() string {
	return "call"
}
