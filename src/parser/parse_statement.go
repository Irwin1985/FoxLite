package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStmt()
	default:
		if p.match(token.Ident) && p.peekToken.Type == token.Assign {
			return p.parseVarStmt()
		}
		return p.parseExpressionStmt()
	}
}
