package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/token"
)

func (p *Parser) parseClassStmt() ast.Statement {
	stmt := &ast.Class{
		Token:      p.curToken,
		Properties: map[string]ast.Expression{},
		Methods:    map[string]*ast.FunctionLiteral{},
	}
	p.nextToken() // skip 'class' token
	if !p.match(token.Ident) {
		p.newError("invalid class name")
	}
	stmt.Name = p.curToken.Literal
	p.nextToken() // skip 'ident' token

	p.expect(token.NewLine, "")

	stmt.Properties = p.parseClassProperties()
	stmt.Methods = p.parseMethods()

	return stmt
}

func (p *Parser) parseClassProperties() map[string]ast.Expression {
	prop := make(map[string]ast.Expression)

	for !p.eof() && p.isProperty() {
		key := p.curToken.Literal
		p.nextToken() // skip 'Ident' token
		p.nextToken() // skip '=' token
		val := p.parseExpression(lowest)
		prop[key] = val
		p.expect(token.NewLine, "")
	}

	return prop
}

func (p *Parser) parseMethods() map[string]*ast.FunctionLiteral {
	var fns = make(map[string]*ast.FunctionLiteral)
	for !p.eof() && p.match(token.Function) {
		val := p.parseFunctionLiteral().(*ast.FunctionLiteral)
		key := val.Name.Value.(string)
		fns[key] = val
	}
	return fns
}

func (p *Parser) isProperty() bool {
	return p.curToken.Type == token.Ident && p.peekToken.Type == token.Assign
}
