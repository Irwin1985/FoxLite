package parser

import "FoxLite/src/ast"

func (p *Parser) parseWhileStmt() ast.Statement {
	stmt := &ast.While{
		Token: p.curToken,
	}

	p.nextToken() // skip 'While' token
	stmt.Condition = p.parseExpression(lowest)
	stmt.Body = p.parseBlockStmt()

	return stmt
}
