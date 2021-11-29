package parser

import "FoxLite/src/ast"

func (p *Parser) parseInfixExp(leftExp ast.Expression) ast.Expression {
	exp := &ast.InfixExp{
		Token: p.curToken,
		Left:  leftExp,
		Op:    p.curToken.Type,
	}
	precedence := p.curPrecedence() // guardamos el orden de precedencia del operador
	p.nextToken()                   // avanza el operador
	exp.Right = p.parseExpression(precedence)

	return exp
}
