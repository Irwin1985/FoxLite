package parser

import "FoxLite/src/ast"

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{
		Token: p.curToken,
	}
	p.nextToken() // skip 'return' keyword
	stmt.Value = p.parseExpression(lowest)
	return stmt
}
