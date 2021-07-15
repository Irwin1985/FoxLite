package lexer

import (
	"FoxLite/lang/token"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	s   *Stream
	lt  token.TokenType
	c   rune
	ln  int
	col int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		s:  NewStream(input),
		lt: token.EOF,
	}
	l.eat() // prime the first character
	return l
}

// create a token and keep track of last generated token
func (l *Lexer) newToken(t token.TokenType, v interface{}) token.Token {
	tok := token.Token{Type: t, Lexeme: v, Ln: l.ln, Col: l.col}
	l.lt = t
	return tok
}

func isWS(c rune) bool {
	return c == rune(' ') || c == rune('\t') || c == rune('\r')
}

func isAZ(c rune) bool {
	return unicode.IsLetter(c) || c == rune('_')
}

func isID(c rune) bool {
	return isAZ(c) || unicode.IsNumber(c)
}

func (l *Lexer) eat() {
	l.c = l.s.Read()
}

func (l *Lexer) isEOF() bool {
	return l.c == rune(0)
}

func (l *Lexer) ws() {
	for !l.isEOF() && isWS(l.c) {
		l.eat()
	}
}

func (l *Lexer) getNum() string {
	v := ""
	for !l.isEOF() && unicode.IsNumber(l.c) {
		v += string(l.c)
		l.eat()
	}
	return v
}

func (l *Lexer) num() token.Token {
	lex := string(l.c)
	l.eat()
	lex += l.getNum()
	if l.c == '.' && unicode.IsDigit(l.s.Peek()) {
		lex += "."
		l.eat()
		lex += l.getNum()
	}

	v, err := strconv.ParseFloat(lex, 64)
	if err != nil {
		fmt.Printf("could not parse the following lexeme into number token: %v", lex)
		os.Exit(1)
	}

	return l.newToken(token.NUMBER, float32(v))
}

func (l *Lexer) str() token.Token {
	strEnd := l.c
	l.eat()
	var lex string
	for !l.isEOF() && l.c != strEnd {
		lex += string(l.c)
		l.eat()
	}
	l.eat()
	return l.newToken(token.STRING, lex)
}

func (l *Lexer) ident() token.Token {
	var v string
	for !l.isEOF() && isID(l.c) {
		v += string(l.c)
		l.eat()
	}
	v = strings.ToLower(v)
	return l.newToken(token.LookupIdent(v), v)
}

func (l *Lexer) punc(t token.TokenType) token.Token {
	c1 := string(l.c)
	l.eat()
	twoc := c1 + string(l.c)

	if tok, ok := token.Special(twoc); ok {
		l.eat()
		return l.newToken(tok, twoc)
	}
	return l.newToken(t, c1)

}

func isBoolStr(c rune) bool {
	return c == 't' || c == 'T' || c == 'f' || c == 'F'
}

func (l *Lexer) dot() token.Token {
	// we need to recognize this pattern '.' 't|T|f|F' '.'
	l.eat() // skip the first '.'
	if isBoolStr(l.c) && l.s.Peek() == '.' {
		// defenitely is a boolean token
		isFalse := l.c == 'f' || l.c == 'F'
		l.eat() // skip 't' | 'T' | 'f' | 'F'
		l.eat() // skip '.'
		tok := token.TRUE
		v := true
		if isFalse {
			tok = token.FALSE
			v = false
		}
		return l.newToken(tok, v)
	}
	return l.newToken(token.DOT, ".")
}

func (l *Lexer) NextToken() token.Token {
	for !l.isEOF() {
		l.ln = l.s.Line
		l.col = l.s.Col
		if l.c == rune(';') {
			l.eat()
			l.lt = token.SEMICOLON
		}
		if l.c == rune('\n') {
			if l.lt == token.EOF || l.lt == token.NEWLINE || l.lt == token.SEMICOLON {
				l.eat()
				continue
			} else {
				return l.newToken(token.NEWLINE, "\\n")
			}
		}
		if isWS(l.c) {
			l.ws()
			continue
		}
		if unicode.IsNumber(l.c) {
			return l.num()
		}
		if isAZ(l.c) {
			return l.ident()
		}
		if l.c == rune('\'') || l.c == rune('"') {
			return l.str()
		}
		if l.c == rune('.') {
			return l.dot()
		}
		if tok, ok := token.Special(string(l.c)); ok {
			return l.punc(tok)
		}
		fmt.Printf("Unknown character at [%d:%d]\n", l.ln, l.col)
		os.Exit(1)
	}
	return token.Token{Type: token.EOF, Lexeme: ""}
}
