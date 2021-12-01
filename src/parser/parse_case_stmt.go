package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseDoCaseStmt() ast.Statement {
	stmt := &ast.DoCaseStmt{
		Token:      p.curToken,
		Statements: []*ast.CaseBranch{},
	}
	p.nextToken()               // skip 'Do' token
	p.nextToken()               // skip 'Case' token
	p.expect(token.NewLine, "") // skip newline token

	for !p.eof() && p.match(token.Case) {
		p.nextToken() // skip 'Case' token
		caseBranch := &ast.CaseBranch{
			Condition: p.parseExpression(lowest),
			Body:      p.parseBlockStmt(),
		}
		stmt.Statements = append(stmt.Statements, caseBranch)
	}

	if p.match(token.Otherwise) {
		p.nextToken() // skip 'Otherwise' token
		stmt.Alternative = p.parseBlockStmt()
		p.expect(token.EndCase, "")
	}

	return stmt
}
