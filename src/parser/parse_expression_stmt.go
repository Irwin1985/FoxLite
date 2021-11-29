package parser

import "FoxLite/src/ast"

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{
		Token: p.curToken,
	}
	exp := p.parseExpression(lowest)
	if exp != nil {
		stmt.Expression = exp
	}
	return stmt
}
