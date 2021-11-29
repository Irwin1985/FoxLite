package ast

import (
	"FoxLite/src/token"
	"fmt"
)

type Literal struct {
	Token token.Token // puede ser: ident, number, string, true, false
	Value interface{}
}

func (l *Literal) expressionNode() {}
func (l *Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}
