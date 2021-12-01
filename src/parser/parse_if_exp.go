package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseIfExp() ast.Expression {
	exp := &ast.IfExp{
		Token: p.curToken,
	}
	p.nextToken() // skip 'If' token
	exp.Condition = p.parseExpression(lowest)
	exp.Consequence = p.parseBlockStmt()

	if p.match(token.Else) {
		p.nextToken() // skip 'Else' token
		exp.Alternative = p.parseBlockStmt()
	}

	return exp
}
