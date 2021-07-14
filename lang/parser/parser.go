package parser

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/lexer"
	"FoxLite/lang/token"
	"fmt"
)

/*
* Parser infrastructure
* 1. Precedence order constants
* 2. Token-Precedence association
* 3. Token-Semantic code association
 */

const (
	LOWEST int = iota
	ASSIGNMENT
	LOGIC_OR
	LOGIC_AND
	EQUALITY   // '==' | '!='
	COMPARISON // '<' | '<=' | '>' | '>='
	TERM       // '+' | '-'
	FACTOR     // '*' | '/'
	UNARY      // '!' | '-'
	CALL       // foo()
	INDEX      // foo[bar]
)

var mapPrecedence = map[token.TokenType]int{
	// logical arithmetic
	token.OR:  LOGIC_OR,
	token.AND: LOGIC_AND,
	// assignment
	token.ASSIGN: ASSIGNMENT,
	// equality
	token.EQ:  EQUALITY,
	token.NEQ: EQUALITY,
	// comparison
	token.LT:  COMPARISON,
	token.GT:  COMPARISON,
	token.GEQ: COMPARISON,
	token.LEQ: COMPARISON,
	// term
	token.PLUS:     TERM,
	token.PLUS_EQ:  TERM,
	token.MINUS:    TERM,
	token.MINUS_EQ: TERM,
	// factor
	token.MUL:    FACTOR,
	token.MUL_EQ: FACTOR,
	token.DIV:    FACTOR,
	token.DIV_EQ: FACTOR,
	// call
	token.LPAREN: CALL,
}

// parsing functions type
type prefixFn func() ast.Expr
type infixFn func(ast.Expr) ast.Expr

type Parser struct {
	l           *lexer.Lexer
	curToken    token.Token
	peekToken   token.Token
	prevToken   token.Token
	mapPrefixFn map[token.TokenType]prefixFn
	mapInfixFn  map[token.TokenType]infixFn
	errors      []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:           l,
		mapPrefixFn: make(map[token.TokenType]prefixFn),
		mapInfixFn:  make(map[token.TokenType]infixFn),
		errors:      []string{},
	}

	p.advance()
	p.advance()

	p.regSemanticCode()

	return p
}

func (p *Parser) regSemanticCode() {
	// Semantic code for prefix tokens
	p.regPrefixFn(token.IDENT, p.parseLiteralExpr)
	p.regPrefixFn(token.NUMBER, p.parseLiteralExpr)
	p.regPrefixFn(token.STRING, p.parseLiteralExpr)
	p.regPrefixFn(token.TRUE, p.parseLiteralExpr)
	p.regPrefixFn(token.FALSE, p.parseLiteralExpr)
	p.regPrefixFn(token.NULL, p.parseLiteralExpr)
	p.regPrefixFn(token.MINUS, p.parseUnaryExpr)
	p.regPrefixFn(token.BANG, p.parseUnaryExpr)
	p.regPrefixFn(token.LPAREN, p.parseGroupedExpr)

	// Semantic code for infix tokens
	p.regInfixFn(token.PLUS, p.parseBinaryExpr)
	p.regInfixFn(token.PLUS_EQ, p.parseBinaryExpr)
	p.regInfixFn(token.MINUS, p.parseBinaryExpr)
	p.regInfixFn(token.MINUS_EQ, p.parseBinaryExpr)
	p.regInfixFn(token.MUL, p.parseBinaryExpr)
	p.regInfixFn(token.MUL_EQ, p.parseBinaryExpr)
	p.regInfixFn(token.DIV, p.parseBinaryExpr)
	p.regInfixFn(token.DIV_EQ, p.parseBinaryExpr)

	p.regInfixFn(token.LT, p.parseBinaryExpr)
	p.regInfixFn(token.GT, p.parseBinaryExpr)
	p.regInfixFn(token.LEQ, p.parseBinaryExpr)
	p.regInfixFn(token.GEQ, p.parseBinaryExpr)

	p.regInfixFn(token.AND, p.parseBinaryExpr)
	p.regInfixFn(token.OR, p.parseBinaryExpr)
	p.regInfixFn(token.ASSIGN, p.parseBinaryExpr)
	p.regInfixFn(token.LPAREN, p.parseCallExpr)
}

func (p *Parser) Parse() []ast.Stmt {
	program := []ast.Stmt{}
	for !p.isAtEnd() {
		s := p.statement()
		if s == nil {
			return nil
		}
		program = append(program, s)
		p.match(token.NEWLINE)
	}
	if !p.curTokenIs(token.EOF) {
		p.newError("expected EOF.")
	}
	return program
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.LOCAL) || p.match(token.PRIVATE) || p.match(token.PUBLIC) {
		return p.varStatement()
	} else if p.match(token.FUNCTION) {
		return p.functionStmt()
	} else if p.match(token.RETURN) {
		return p.returnStmt()
	} else if p.match(token.IF) {
		return p.ifStatement()
	} else {
		return p.expressionStatement()
	}
}

func (p *Parser) functionStmt() ast.Stmt {
	stmt := &ast.FunctionStmt{}
	stmt.Parameters = []ast.LiteralExpr{}

	p.expect(token.IDENT, "expected function name.")
	stmt.Name = p.prevToken

	if p.match(token.LPAREN) {
		if !p.match(token.RPAREN) {
			p.match(token.IDENT)
			stmt.Parameters = append(stmt.Parameters, ast.LiteralExpr{Value: p.prevToken})

			for !p.isAtEnd() && p.match(token.COMMA) {
				p.match(token.IDENT)
				stmt.Parameters = append(stmt.Parameters, ast.LiteralExpr{Value: p.prevToken})
			}
			//p.expect(token.RPAREN, "expected ')' after parameters")
			if !p.match(token.RPAREN) {
				p.newError("expected ')' after parameters")
			}
		}
	}

	stmt.Body = p.parseBlockStmt(token.ENDFUNC)

	return stmt
}

func (p *Parser) returnStmt() ast.Stmt {
	stmt := &ast.ReturnStmt{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	stmt.Value = exp

	return stmt
}

func (p *Parser) ifStatement() ast.Stmt {
	stmt := &ast.IfStmt{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	stmt.Condition = exp
	p.match(token.THEN)
	stmt.Consequence = p.parseBlockStmt(token.ELSE, token.ENDIF)

	if p.prevToken.Type == token.ELSE {
		stmt.Alternative = p.parseBlockStmt(token.ENDIF)
	}
	return stmt
}

func (p *Parser) parseBlockStmt(t ...token.TokenType) *ast.BlockStmt {
	block := &ast.BlockStmt{}
	block.Statements = []ast.Stmt{}
	p.expect(token.NEWLINE, "expected NEWLINE before block")

	for !p.isAtEnd() && !p.match(t...) {
		s := p.statement()
		if s == nil {
			return nil
		}
		block.Statements = append(block.Statements, s)
		p.match(token.NEWLINE)
	}

	return block
}

func (p *Parser) varStatement() ast.Stmt {
	tok := p.prevToken
	if p.match(token.LPAREN) {
		return p.parseGroupedVarStmt()
	} else {
		p.expect(token.IDENT, "expected variable name.")
		name := &ast.LiteralExpr{Value: p.prevToken}
		p.expect(token.ASSIGN, "expected '=' before expression.")
		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}
		if p.curTokenIs(token.COMMA) {
			return p.parseInlineVarStmt(tok, name, value)
		} else {
			stmt := &ast.VarStmt{
				Token: tok,
				Name:  name,
				Value: value,
			}
			return stmt
		}
	}
}

// inlineVarStmt ::= ('LOCAL'|'PRIVATE'|'PUBLIC') identifier '=' varStmt (',' varStmt)*
func (p *Parser) parseInlineVarStmt(tok token.Token, name *ast.LiteralExpr, value ast.Expr) *ast.InlineVarStmt {
	stmt := &ast.InlineVarStmt{}
	stmt.Scope = tok.Type
	stmt.Variables = []ast.Stmt{}

	// add previous variable
	item := &ast.VarStmt{Name: name, Value: value}
	stmt.Variables = append(stmt.Variables, item)

	for !p.isAtEnd() && p.match(token.COMMA) {
		p.expect(token.IDENT, "expected variable name.")

		name = &ast.LiteralExpr{Value: p.prevToken}
		p.expect(token.ASSIGN, "expected '=' before expression.")
		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		item = &ast.VarStmt{Name: name, Value: value}
		stmt.Variables = append(stmt.Variables, item)
	}
	return stmt
}

// groupVarStmt ::= ('LOCAL'|'PRIVATE'|'PUBLIC') '(' (varStmt ',')* ')'
func (p *Parser) parseGroupedVarStmt() *ast.InlineVarStmt {

	stmt := &ast.InlineVarStmt{}
	stmt.Scope = p.prevToken.Type
	stmt.Variables = []ast.Stmt{}

	if p.match(token.RPAREN) {
		p.newError("expected variable declaration.")
		return nil
	}

	p.match(token.NEWLINE) // in case of a NEW_LINE
	p.expect(token.IDENT, "expected variable name.")

	name := &ast.LiteralExpr{Value: p.prevToken}
	p.expect(token.ASSIGN, "expected '=' before expression.")
	value := p.parseExpression(LOWEST)
	if value == nil {
		return nil
	}

	item := &ast.VarStmt{Name: name, Value: value}
	stmt.Variables = append(stmt.Variables, item)

	for !p.isAtEnd() && p.match(token.COMMA) {

		p.match(token.NEWLINE) // in case of a NEW_LINE
		p.expect(token.IDENT, "expected variable name.")

		name = &ast.LiteralExpr{Value: p.prevToken}
		p.expect(token.ASSIGN, "expected '=' before expression.")
		value = p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		item = &ast.VarStmt{Name: name, Value: value}
		stmt.Variables = append(stmt.Variables, item)
	}

	p.match(token.NEWLINE) // in case of a NEW_LINE
	p.expect(token.RPAREN, "expected ')' after expression.")

	return stmt
}

func (p *Parser) expressionStatement() ast.Stmt {
	stmt := &ast.ExprStmt{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	stmt.Expression = exp
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := p.mapPrefixFn[p.curToken.Type]
	if prefix == nil {
		m := "Function argument value, type, or count is invalid."
		p.newError(fmt.Sprintf("%s\nParsing function not found for token: %v", m, token.TokenNames[p.curToken.Type]))
		return nil
	}
	leftExpr := prefix()

	for !p.isAtEnd() && precedence < p.curPrecedence() {
		infix := p.mapInfixFn[p.curToken.Type]
		if infix == nil {
			return leftExpr
		}
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseLiteralExpr() ast.Expr {
	expr := &ast.LiteralExpr{Value: p.curToken}
	p.advance()
	return expr
}

func (p *Parser) parseBinaryExpr(left ast.Expr) ast.Expr {
	expr := &ast.Binary{Left: left, Operator: p.curToken}

	pre := p.curPrecedence()
	p.advance()
	exp := p.parseExpression(pre)
	if expr == nil {
		return nil
	}
	expr.Right = exp

	return expr
}

func (p *Parser) parseCallExpr(left ast.Expr) ast.Expr {
	expr := &ast.CallExpr{Function: left}

	p.advance()
	if !p.match(token.RPAREN) {
		expr.Arguments = p.parseExpressions()
		p.expect(token.RPAREN, "expected ')' after arguments")
	}

	return expr
}

func (p *Parser) parseExpressions() []ast.Expr {
	exprList := []ast.Expr{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	exprList = append(exprList, exp)
	for !p.isAtEnd() && p.match(token.COMMA) {
		exp := p.parseExpression(LOWEST)
		if exp == nil {
			return nil
		}
		exprList = append(exprList, exp)
	}
	return exprList
}

func (p *Parser) parseUnaryExpr() ast.Expr {
	expr := &ast.Unary{Operator: p.curToken}
	p.advance()
	exp := p.parseExpression(UNARY)
	if exp == nil {
		return nil
	}
	expr.Right = exp

	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.advance()
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	p.expect(token.RPAREN, "expected ')' after expression.")

	return exp
}

func (p *Parser) regPrefixFn(t token.TokenType, fn prefixFn) {
	p.mapPrefixFn[t] = fn
}

func (p *Parser) regInfixFn(t token.TokenType, fn infixFn) {
	p.mapInfixFn[t] = fn
}

func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expect(t token.TokenType, msg string) {
	if !p.match(t) {
		msg := fmt.Sprintf("Syntax error at (%d:%d) %s.", p.curToken.Ln, p.curToken.Col, msg)
		p.newError(msg)
	}
}

func (p *Parser) curTokenIs(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.curToken.Type == t {
			p.prevToken = p.curToken
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.curToken.Type == token.EOF
}

func (p *Parser) newError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) curPrecedence() int {
	if pre, ok := mapPrecedence[p.curToken.Type]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}
