package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parsePrintStmt() ast.Statement {
	stmt := &ast.PrintStmt{
		Token:    p.curToken,
		Messages: []ast.Expression{},
	}
	p.nextToken() // skip '?' token
	stmt.Messages = append(stmt.Messages, p.parseExpression(lowest))

	for !p.eof() && p.match(token.Comma) {
		p.nextToken() // skip ',' token
		stmt.Messages = append(stmt.Messages, p.parseExpression(lowest))
	}

	return stmt
}
