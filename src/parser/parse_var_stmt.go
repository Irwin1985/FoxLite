package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
	"fmt"
)

func (p *Parser) parseVarStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{
		Token: p.curToken,
		Scope: 'p', // private
		Type:  'b', // boolean
		Value: &ast.Literal{
			Token: token.Token{
				Type:    token.False,
				Literal: "false",
			},
			Value: false},
	}

	if p.match(token.Local, token.Private, token.Public) {
		switch p.curToken.Type {
		case token.Local:
			stmt.Scope = 'l'
		case token.Private:
			stmt.Scope = 'p'
		case token.Public:
			stmt.Scope = 'g'
		default:
			stmt.Scope = 'p' // private
		}
		p.nextToken() // skip variable scope (local, private, public)
	}
	// nombre de la variable
	if !p.match(token.Ident) {
		p.newError(fmt.Sprintf("unexpected token `%s` for variable name", p.curToken.Literal))
	}
	stmt.Name = p.curToken.Literal
	p.nextToken() // skip variable name

	// opcionalmente se puede asignar un valor
	if p.match(token.Assign) {
		p.nextToken() // skip '=' token
		stmt.Value = p.parseExpression(lowest)
	}

	return stmt
}
