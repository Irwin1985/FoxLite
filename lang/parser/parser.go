package parser

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/lexer"
	"FoxLite/lang/token"
	"fmt"
	"os"
	"strconv"
)

// precedence order
const (
	LOWEST int = iota
	LOGIC_OR
	LOGIC_AND
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	PREFIX
	CALL
	INDEX
)

// precedence table
var precedence = map[token.TokenType]int{
	// logic arithmetic
	token.OR:  LOGIC_OR,
	token.AND: LOGIC_AND,
	// equality
	token.EQ:  EQUALITY,
	token.NEQ: EQUALITY,
	// comparison
	token.LT:  COMPARISON,
	token.LEQ: COMPARISON,
	token.GT:  COMPARISON,
	token.GEQ: COMPARISON,
	// term
	token.PLUS:  TERM,
	token.MINUS: TERM,
	// factor
	token.MUL: FACTOR,
	token.DIV: FACTOR,
	// call
	token.LPAREN: CALL,
	// index
	token.LBRACKET: INDEX,
}

type Parser struct {
	l             *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	prefixParseFn map[token.TokenType]func() ast.Expression
	infixParseFn  map[token.TokenType]func(ast.Expression) ast.Expression
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	// prefix function association
	p.prefixParseFn = make(map[token.TokenType]func() ast.Expression)
	p.prefixParseFn[token.IDENT] = p.parseIdentifier
	p.prefixParseFn[token.INT] = p.parseIntegerLiteral
	p.prefixParseFn[token.STRING] = p.parseStringLiteral
	p.prefixParseFn[token.TRUE] = p.parseBooleanLiteral
	p.prefixParseFn[token.FALSE] = p.parseBooleanLiteral
	p.prefixParseFn[token.NULL] = p.parseNilLiteral
	p.prefixParseFn[token.IF] = p.parseIfExpression

	// infix function association
	p.infixParseFn = make(map[token.TokenType]func(ast.Expression) ast.Expression)
	p.infixParseFn[token.PLUS] = p.parseInfixExpression
	p.infixParseFn[token.MINUS] = p.parseInfixExpression
	p.infixParseFn[token.MUL] = p.parseInfixExpression
	p.infixParseFn[token.DIV] = p.parseInfixExpression
	p.infixParseFn[token.LT] = p.parseInfixExpression
	p.infixParseFn[token.LEQ] = p.parseInfixExpression
	p.infixParseFn[token.GT] = p.parseInfixExpression
	p.infixParseFn[token.GEQ] = p.parseInfixExpression
	p.infixParseFn[token.EQ] = p.parseInfixExpression
	p.infixParseFn[token.NEQ] = p.parseInfixExpression

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expect(t token.TokenType) {
	if p.curToken.Type == t {
		p.nextToken()
	} else {
		fmt.Printf("expected %d and got %d\n", t, p.curToken.Type)
		os.Exit(1)
	}
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedence[p.curToken.Type]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) curTokenIs(t ...token.TokenType) bool {
	for _, tok := range t {
		if p.curToken.Type == tok {
			return true
		}
	}
	return false
}

// program ::= (stmt)* EOF
func (p *Parser) Program() *ast.Program {
	var program = &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		stmt := p.stmt()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.EOF) {
		fmt.Println("expected end of file")
		os.Exit(1)
	}

	return program
}

// stmt ::= varStmt | returnStmt | exprStmt
func (p *Parser) stmt() ast.Statement {
	switch p.curToken.Type {
	case token.LOCAL:
		return p.varStmt()
	case token.RETURN:
		return p.returnStmt()
	default:
		if p.curTokenIs(token.IDENT) && p.peekToken.Type == token.BINDING {
			return p.bindingStmt()
		}
		return p.exprStmt()
	}
}

// varStmt ::= 'var' ident '=' expression
func (p *Parser) varStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{Token: p.curToken}
	p.nextToken()

	if !p.curTokenIs(token.IDENT) {
		fmt.Println("Identifier expected")
		os.Exit(1)
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Lexeme}
	p.nextToken()

	p.expect(token.ASSIGN)

	stmt.Value = p.expression(LOWEST)

	return stmt
}

// bindingStmt ::= varStmt
func (p *Parser) bindingStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{Token: p.curToken}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Lexeme}
	p.nextToken() // ident

	p.nextToken() // ':='

	stmt.Value = p.expression(LOWEST)

	return stmt
}

// returnStmt ::= 'return' expression
func (p *Parser) returnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.expression(LOWEST)

	return stmt
}

// exprStmt
func (p *Parser) exprStmt() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.expression(LOWEST)
	return stmt
}

// expression
func (p *Parser) expression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFn[p.curToken.Type]
	if prefixFn == nil {
		return nil
	}
	left := prefixFn()
	for precedence < p.curPrecedence() {
		infixFn := p.infixParseFn[p.curToken.Type]
		if infixFn == nil {
			return left
		}
		left = infixFn(left)
	}
	return left
}

// parseIdentifier
func (p *Parser) parseIdentifier() ast.Expression {
	exp := &ast.Identifier{Token: p.curToken, Value: p.curToken.Lexeme}
	p.nextToken()
	return exp
}

// parseIntegerLiteral
func (p *Parser) parseIntegerLiteral() ast.Expression {
	exp := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Lexeme, 0, 64)
	if err != nil {
		fmt.Printf("could not parse %q as integer", p.curToken.Lexeme)
		os.Exit(1)
	}
	exp.Value = value
	p.nextToken()
	return exp
}

// parseStringLiteral
func (p *Parser) parseStringLiteral() ast.Expression {
	exp := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Lexeme}
	p.nextToken()
	return exp
}

// parseBooleanLiteral
func (p *Parser) parseBooleanLiteral() ast.Expression {
	exp := &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	p.nextToken()
	return exp
}

// parseNilLiteral
func (p *Parser) parseNilLiteral() ast.Expression {
	exp := &ast.NilLiteral{Token: p.curToken}
	p.nextToken()
	return exp
}

// parseIfExpression
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	p.nextToken()

	exp.Condition = p.expression(LOWEST)

	exp.Consequence = p.parseBlockStmt(token.ELSE, token.ENDIF)

	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		exp.Alternative = p.parseBlockStmt(token.ENDIF)
	}
	p.expect(token.ENDIF)

	return exp
}

// parseBlockStmt
func (p *Parser) parseBlockStmt(t ...token.TokenType) *ast.BlockStmt {
	block := &ast.BlockStmt{}
	block.Statements = []ast.Statement{}
	p.expect(token.NEWLINE)
	for !p.curTokenIs(t...) {
		stmt := p.stmt()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}
	}

	return block
}

// parseInfixExpression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.curToken, Op: p.curToken.Lexeme, Left: left}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.expression(precedence)

	return exp
}
