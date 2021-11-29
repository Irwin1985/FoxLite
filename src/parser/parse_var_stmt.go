package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
	"strings"
)

func (p *Parser) parseVarStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{
		Token: p.curToken,
		Name:  p.curToken.Literal,
		Scope: 'p', // private
		Type:  'b', // boolean
	}
	// nombre de la variable
	varName := stmt.Name
	p.nextToken() // avanzamos el token ident

	// detectar el ámbito de la variable
	varScope := varName[0:1]
	if strings.Contains("lpg", varScope) {
		stmt.Scope = varScope[0]
		// Detectar el tipo de dato
		varType := varName[1:2]
		if strings.Contains("cnbod", varType) {
			stmt.Type = varType[0]
		}
	}

	p.expect(token.Assign, "expecting `=` (e.g. `lcName = 'Jhon'`)")
	stmt.Value = p.parseExpression(lowest)

	return stmt
}
