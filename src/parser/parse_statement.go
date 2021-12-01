package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStmt()
	case token.CloseQM:
		return p.parsePrintStmt()
	default:
		if p.match(token.Ident) && p.peekToken.Type == token.Assign {
			return p.parseVarStmt()
		}
		if p.match(token.Do) && p.peekToken.Type == token.Case {
			return p.parseDoCaseStmt()
		}
		return p.parseExpressionStmt()
	}
}
