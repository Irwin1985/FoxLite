package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
	"strconv"
)

func (p *Parser) parseLiteral() ast.Expression {
	exp := &ast.Literal{
		Token: p.curToken,
	}
	switch p.curToken.Type {
	case token.Number:
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.newError("float conversion failure.")
		}
		exp.Value = val
	case token.String, token.Ident:
		exp.Value = p.curToken.Literal
	case token.Null:
		exp.Value = nil
	case token.True:
		exp.Value = true
	case token.False:
		exp.Value = false
	}
	p.nextToken()
	return exp
}
