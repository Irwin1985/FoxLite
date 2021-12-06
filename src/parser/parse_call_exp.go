package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseCallExp(caller ast.Expression) ast.Expression {
	exp := &ast.CallExp{ // foo(x, y)
		Token:  p.curToken,
		Caller: caller,
		Args:   []ast.Expression{},
	}
	p.nextToken() // skip '(' token

	if !p.match(token.Rparen) {
		exp.Args = append(exp.Args, p.parseExpression(lowest))

		for !p.eof() && p.match(token.Comma) {
			p.nextToken() // skip ',' token
			exp.Args = append(exp.Args, p.parseExpression(lowest))
		}
	}
	p.expect(token.Rparen, "")

	return exp
}
