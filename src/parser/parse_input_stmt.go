package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseInputStmt() ast.Statement {
	stmt := &ast.Input{
		Token: p.curToken,
	}
	p.nextToken() // skip 'OpenQM' token
	stmt.Message = p.parseExpression(lowest)
	p.expect(token.Comma, "")
	stmt.Output = p.parseLiteral().(*ast.Literal)

	return stmt
}
