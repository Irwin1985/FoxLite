package ast

import (
	"FoxLite/src/token"
	"bytes"
	"strings"
)

type FunctionLiteral struct {
	Token      token.Token
	Name       *Literal
	Parameters []*Literal
	Body       *BlockStmt
}

func (f *FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("Func")
	out.WriteString(f.Name.String())
	out.WriteString("(")

	if len(f.Parameters) > 0 {
		var params []string
		for _, param := range f.Parameters {
			params = append(params, param.String())
		}
		out.WriteString(strings.Join(params, ", "))
	}
	out.WriteString(")")
	out.WriteString(f.Body.String())
	return out.String()
}
