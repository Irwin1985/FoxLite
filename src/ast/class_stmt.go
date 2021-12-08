package ast

import "FoxLite/src/token"

type Class struct {
	Token      token.Token
	Name       string
	Properties map[string]Expression
	Methods    map[string]*FunctionLiteral
}

func (c *Class) statementNode() {}
func (c *Class) String() string {
	return "class"
}
