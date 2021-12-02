package ast

import "FoxLite/src/token"

type Exit struct {
	Token token.Token
}

func (e *Exit) statementNode() {}
func (e *Exit) String() string {
	return "exit"
}
