package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
	"fmt"
)

func (p *Parser) parseFunctionLiteral() ast.Statement {
	exp := &ast.FunctionLiteral{
		Token:      p.curToken,
		Parameters: []*ast.Literal{},
	}
	p.nextToken() // skip 'Func' token
	exp.Name = p.parseLiteral().(*ast.Literal)
	p.expect(token.Lparen, fmt.Sprintf("unexpected token `%s`, expecting `(`", p.curToken.Literal))

	if !p.match(token.Rparen) { // hay parámetros definidos?
		exp.Parameters = append(exp.Parameters, p.parseLiteral().(*ast.Literal))
		for !p.eof() && p.match(token.Comma) {
			p.nextToken() // skip ',' token
			exp.Parameters = append(exp.Parameters, p.parseLiteral().(*ast.Literal))
		}
	}
	p.expect(token.Rparen, "expecting `)`")

	exp.Body = p.parseBlockStmt()

	return exp
}
