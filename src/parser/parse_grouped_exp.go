package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseGroupedExp() ast.Expression {
	p.nextToken() // skip '('
	exp := p.parseExpression(lowest)
	p.expect(token.Rparen, "expected `)` after grouped expresion")
	return exp
}
