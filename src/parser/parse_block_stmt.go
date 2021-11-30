package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{
		Statements: []ast.Statement{},
	}
	p.expect(token.NewLine, "cannot call a function that does not have a body")
	col := p.curToken.Col
	// Parseamos todos los tokens que
	// coincidan con la columna 'col'
	for !p.eof() && p.curToken.Col == col {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.match(token.NewLine) {
			p.nextToken()
		}
	}

	return block
}
