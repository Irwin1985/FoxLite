package ast

import "FoxLite/src/token"

type Input struct {
	Token   token.Token
	Message Expression
	Output  *Literal
}

func (i *Input) statementNode() {}
func (i *Input) String() string {
	return "input"
}
