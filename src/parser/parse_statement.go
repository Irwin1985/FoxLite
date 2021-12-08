package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Return:
		return p.parseReturnStmt()
	case token.OpenQM:
		return p.parseInputStmt()
	case token.CloseQM:
		return p.parsePrintStmt()
	case token.While:
		return p.parseWhileStmt()
	case token.Loop:
		stmt := &ast.Loop{Token: p.curToken}
		p.nextToken()
		return stmt
	case token.Exit:
		stmt := &ast.Exit{Token: p.curToken}
		p.nextToken()
		return stmt
	case token.Class:
		return p.parseClassStmt()
	case token.Function:
		return p.parseFunctionLiteral()
	case token.If:
		return p.parseIfStmt()
	default:
		if p.isVarStmt() {
			return p.parseVarStmt()
		}
		if p.match(token.Do) && p.peekToken.Type == token.Case {
			return p.parseDoCaseStmt()
		}
		return p.parseExpressionStmt()
	}
}

func (p *Parser) isVarStmt() bool {
	if p.match(token.Ident) && p.peekToken.Type == token.Assign {
		return true
	}
	return p.match(token.Local, token.Private, token.Public)
}
