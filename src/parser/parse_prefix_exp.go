package parser

import "FoxLite/src/ast"

func (p *Parser) parsePrefixExp() ast.Expression {
	exp := &ast.PrefixExp{
		Token: p.curToken,
		Op:    p.curToken.Type,
	}
	p.nextToken() // avanza el token prefix (!, -)
	exp.Right = p.parseExpression(index)

	return exp
}
