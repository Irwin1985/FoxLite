package parser

import (
	"FoxLite/src/ast"
	"fmt"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFns := p.prefixParseFns[p.curToken.Type]
	if prefixFns == nil {
		p.newError(fmt.Sprintf("unexpected token `%s`", p.curToken.Literal))
		p.recovery()
		return nil
	}
	// parse left right expresion
	leftExp := prefixFns()

	for precedence < p.curPrecedence() {
		infixFns := p.infixParseFns[p.curToken.Type]
		if infixFns == nil {
			return leftExp
		}
		leftExp = infixFns(leftExp)
	}
	return leftExp
}
