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
	l.consume() // prime the first character
	return l
}

// create a token and keep track of last generated token
func (l *Lexer) newToken(t token.TokenType, v interface{}) token.Token {
	tok := token.Token{Type: t, Lexeme: v, Ln: l.ln, Col: l.col}
	l.lt = t
	return tok
}

func isSpace(c rune) bool {
	return c == rune(' ') || c == rune('\t') || c == rune('\r')
}

func isLetter(c rune) bool {
	return unicode.IsLetter(c) || c == rune('_')
}

func isIdent(c rune) bool {
	return isLetter(c) || unicode.IsNumber(c)
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

func (l *Lexer) justNum() string {
	lex := ""
	for !l.isAtEnd() && unicode.IsNumber(l.c) {
		lex += string(l.c)
		l.consume()
	}
	return lex
}

func (l *Lexer) num() token.Token {
	lex := string(l.c)
	l.consume()
	lex += l.justNum()
	if l.c == '.' && unicode.IsDigit(l.s.Peek()) {
		lex += "."
		l.consume()
		lex += l.justNum()
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

func isBoolLetter(c rune) bool {
	return c == 't' || c == 'T' || c == 'f' || c == 'F'
}

func (l *Lexer) dot() token.Token {
	// we need to recognize this pattern '.' 't|T|f|F' '.'
	l.consume() // skip the first '.'
	if isBoolLetter(l.c) && l.s.Peek() == '.' {
		// defenitely is a boolean token
		isFalse := l.c == 'f' || l.c == 'F'
		l.consume() // skip 't' | 'T' | 'f' | 'F'
		l.consume() // skip '.'
		tok := token.TRUE
		lex := true
		if isFalse {
			tok = token.FALSE
			lex = false
		}
		return l.newToken(tok, lex)
	}
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
		if unicode.IsNumber(l.c) {
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
