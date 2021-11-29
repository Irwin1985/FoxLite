package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type InfixExp struct {
	Token token.Token
	Left  Expression
	Op    token.TokenType
	Right Expression
}

func (i *InfixExp) expressionNode() {}
func (i *InfixExp) String() string {
	return fmt.Sprintf("%v %v %v", i.Left.String(), i.Op, i.Right.String())
}
