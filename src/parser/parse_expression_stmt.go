package parser

import "FoxLite/src/ast"

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{}
	exp := p.parseExpression(lowest)
	if exp != nil {
		stmt.Expression = exp
	}
	return stmt
}
