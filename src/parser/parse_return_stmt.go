package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{
		Token: p.curToken,
	}
	p.nextToken() // skip 'return' keyword
	if !p.match(token.NewLine) {
		stmt.Value = p.parseExpression(lowest)
	} else {
		stmt.Value = nil
	}
	return stmt
}
