package lexer

import (
	"FoxLite/lang/token"
	"fmt"
	"os"
	"strings"
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
	l.consume() // prime the first character
	return l
}

// create a token and keep track of last generated token
func (l *Lexer) newToken(t token.TokenType, v string) token.Token {
	tok := token.Token{Type: t, Lexeme: v, Line: l.ln, Col: l.col}
	l.lt = t
	return tok
}

func isSpace(c rune) bool {
	return c == rune(' ') || c == rune('\t') || c == rune('\r')
}

func isLetter(c rune) bool {
	return rune('a') <= c && c <= rune('z') || rune('A') <= c && c <= rune('Z') || c == rune('_')
}

func isIdent(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func isDigit(c rune) bool {
	return rune('0') <= c && c <= rune('9')
}

func (l *Lexer) consume() {
	l.c = l.s.Read()
}

func (l *Lexer) isAtEnd() bool {
	return l.c == rune(0)
}

func (l *Lexer) ws() {
	for !l.isAtEnd() && isSpace(l.c) {
		l.consume()
	}
}

func (l *Lexer) num() token.Token {
	var lex string
	for !l.isAtEnd() && isDigit(l.c) {
		lex += string(l.c)
		l.consume()
	}
	return l.newToken(token.INT, lex)
}

func (l *Lexer) str() token.Token {
	strEnd := l.c
	l.consume()
	var lex string
	for !l.isAtEnd() && l.c != strEnd {
		lex += string(l.c)
		l.consume()
	}
	l.consume()
	return l.newToken(token.STRING, lex)
}

func (l *Lexer) ident() token.Token {
	var lex string
	for !l.isAtEnd() && isIdent(l.c) {
		lex += string(l.c)
		l.consume()
	}
	lex = strings.ToLower(lex)
	return l.newToken(token.LookupIdent(lex), lex)
}

func (l *Lexer) special(t token.TokenType) token.Token {
	oldc := string(l.c)
	l.consume()
	twoc := oldc + string(l.c)

	if tok, ok := token.Special(twoc); ok {
		l.consume()
		return l.newToken(tok, twoc)
	}
	return l.newToken(t, oldc)

}

func (l *Lexer) dot() token.Token {
	c := l.s.Peek()
	if c == 't' || c == 'T' || c == 'f' || c == 'F' {
		isFalse := c == 'f' || c == 'F'
		l.consume() // skip '.'
		if l.s.Peek() == '.' {
			l.consume() // skip 't' | 'T'
			l.consume() // skip '.'
			tok := token.TRUE
			lex := ".T."
			if isFalse {
				tok = token.FALSE
				lex = ".F."
			}
			return l.newToken(tok, lex)
		}
		return l.newToken(token.DOT, ".")
	}
	l.consume() // skip '.'
	return l.newToken(token.DOT, ".")
}

func (l *Lexer) NextToken() token.Token {
	for !l.isAtEnd() {
		l.ln = l.s.Line
		l.col = l.s.Col
		if l.c == rune(';') {
			l.consume()
			l.lt = token.SEMICOLON
		}
		if l.c == rune('\n') {
			if l.lt == token.EOF || l.lt == token.NEWLINE || l.lt == token.SEMICOLON {
				l.consume()
				continue
			} else {
				return l.newToken(token.NEWLINE, "\\n")
			}
		}
		if isSpace(l.c) {
			l.ws()
			continue
		}
		if isDigit(l.c) {
			return l.num()
		}
		if isLetter(l.c) {
			return l.ident()
		}
		if l.c == rune('\'') || l.c == rune('"') {
			return l.str()
		}
		if l.c == rune('.') {
			return l.dot()
		}
		if tok, ok := token.Special(string(l.c)); ok {
			return l.special(tok)
		}
		fmt.Printf("Unknown character at [%d:%d]\n", l.ln, l.col)
		os.Exit(1)
	}
	return token.Token{Type: token.EOF, Lexeme: ""}
}
