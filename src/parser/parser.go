package parser

import (
	"FoxLite/src/ast"
	"FoxLite/src/lexer"
	"FoxLite/src/token"
	"fmt"
)

// Constantes con las precedencias
const (
	lowest = iota
	assignment
	logicOr
	logicAnd
	equality
	comparison
	term
	factor
	call
	index
	dot
)

// Tabla de precedencias
var precedenceTable = map[token.TokenType]int{
	token.Assign:    assignment,
	token.Or:        logicOr,
	token.And:       logicAnd,
	token.Equal:     equality,
	token.NotEq:     equality,
	token.Less:      comparison,
	token.LessEq:    comparison,
	token.Greater:   comparison,
	token.GreaterEq: comparison,
	token.Plus:      term,
	token.Minus:     term,
	token.Mul:       factor,
	token.Div:       factor,
	token.Mod:       factor,
	token.Pow:       factor,
	token.Lparen:    call,
	token.Lbracket:  index,
	token.Dot:       dot,
}

type prefixFns = func() ast.Expression
type infixFns = func(left ast.Expression) ast.Expression

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	// Diccionarios para las funciones
	prefixParseFns map[token.TokenType]prefixFns
	infixParseFns  map[token.TokenType]infixFns
	// Informe de errores
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: map[token.TokenType]prefixFns{},
		infixParseFns:  map[token.TokenType]infixFns{},
		errors:         []string{},
	}
	p.registerPrefixFns()
	p.registerInfixFns()
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Program {
	for !p.eof() && p.match(token.NewLine) {
		p.nextToken()
	}
	return p.parseProgram()
}

func (p *Parser) registerPrefixFns() {
	// Constantes literales
	p.prefixParseFns[token.Number] = p.parseLiteral // 123, 45.6
	p.prefixParseFns[token.String] = p.parseLiteral // "foo", 'bar', `xyz`
	p.prefixParseFns[token.True] = p.parseLiteral   // True
	p.prefixParseFns[token.False] = p.parseLiteral  // False
	p.prefixParseFns[token.Null] = p.parseLiteral   // Null
	p.prefixParseFns[token.Ident] = p.parseLiteral  // foo, bar
	// Expresiones agrupadas
	p.prefixParseFns[token.Lparen] = p.parseGroupedExp // (1 + 2) * (3 + 4)
	// Expresiones unarias
	p.prefixParseFns[token.Minus] = p.parsePrefixExp // -5, -foo()
}

func (p *Parser) registerInfixFns() {
	// Operadores aritméticos
	p.infixParseFns[token.Plus] = p.parseInfixExp  // 1 + 2
	p.infixParseFns[token.Minus] = p.parseInfixExp // 1 - 2
	p.infixParseFns[token.Mul] = p.parseInfixExp   // 1 * 3
	p.infixParseFns[token.Div] = p.parseInfixExp   // 1 / 3
	p.infixParseFns[token.Mod] = p.parseInfixExp   // 1 % 3
	p.infixParseFns[token.Pow] = p.parseInfixExp   // 1 ^ 3
	// Operadores lógicos
	p.infixParseFns[token.Or] = p.parseInfixExp  // True or False
	p.infixParseFns[token.And] = p.parseInfixExp // True and False
	// Operadores relacionales
	p.infixParseFns[token.Less] = p.parseInfixExp      // 1 < 2
	p.infixParseFns[token.LessEq] = p.parseInfixExp    // 1 <= 2
	p.infixParseFns[token.Greater] = p.parseInfixExp   // 1 > 2
	p.infixParseFns[token.GreaterEq] = p.parseInfixExp // 1 >= 2
	p.infixParseFns[token.Equal] = p.parseInfixExp     // 1 == 2
	p.infixParseFns[token.NotEq] = p.parseInfixExp     // 1 != 2
	// Operador de resolución de nombres
	p.infixParseFns[token.Dot] = p.parseInfixExp // foo.bar
	// Asignaciones
	p.infixParseFns[token.Assign] = p.parseInfixExp // foo = bar | foo.bar = 20
	// llamadas a funciones
	p.infixParseFns[token.Lparen] = p.parseCallExp // foo()
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedenceTable[p.curToken.Type]; ok {
		return pre
	}
	return lowest
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) expect(t token.TokenType, msg string) {
	if p.curToken.Type != t {
		text := p.l.GetErrorFormat(&p.curToken)
		out := msg
		if len(out) == 0 {
			unexpected := token.GetTokenStr(p.curToken.Type)
			expected := token.GetTokenStr(t)
			out = fmt.Sprintf("unexpected token `%s`, expecting `%s`", unexpected, expected)
		}
		p.errors = append(p.errors, fmt.Sprintf("%s %s\n", text, out))
	} else {
		p.nextToken() // advance the matched token
	}
}

func (p *Parser) peek(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) eof() bool {
	return p.curToken.Type == token.Eof
}

func (p *Parser) newError(msg string) {
	text := p.l.GetErrorFormat(&p.curToken)
	p.errors = append(p.errors, fmt.Sprintf("%s %s\n", text, msg))
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) recovery() {
	// Se recupera tras un error de sintaxis
	for !p.eof() && !p.match(token.NewLine) {
		p.nextToken()
	}
	if p.eof() {
		p.newError("unexpected end of file")
	}
	if p.match(token.NewLine) {
		p.nextToken()
	}
	// Sincronizado
}
