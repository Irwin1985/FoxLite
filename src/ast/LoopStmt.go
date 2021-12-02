package ast

import "FoxLite/src/token"

type Loop struct {
	Token token.Token
}

func (l *Loop) statementNode() {}
func (l *Loop) String() string {
	return "loop"
}
