package lexer

import (
	"FoxLite/src/token"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Lexer struct {
	input     []rune
	scanMode  byte
	fileName  string
	pos       int
	peekPos   int
	ch        rune
	line      int
	col       int
	prevToken token.Token
	symbol    map[string]token.TokenType
	symbols   string
}

func New() *Lexer {
	l := &Lexer{
		symbol:  map[string]token.TokenType{},
		symbols: "+-*/^%",
		line:    1,
		col:     0,
	}
	l.prevToken = l.newToken(token.NewLine, "")

	// Rellenamos los símbolos
	l.symbol["+"] = token.Plus
	l.symbol["+="] = token.PlusEq
	l.symbol["-"] = token.Minus
	l.symbol["-="] = token.MinusEq
	l.symbol["*"] = token.Mul
	l.symbol["*="] = token.MulEq
	l.symbol["/"] = token.Div
	l.symbol["/="] = token.DivEq
	l.symbol["^"] = token.Pow
	l.symbol["%"] = token.Mod

	return l
}

func (l *Lexer) newToken(tok token.TokenType, lit string) token.Token {
	var t = token.Token{
		Type:    tok,
		Literal: lit,
		Col:     l.col,
		Line:    l.line,
	}
	l.prevToken = t
	return t
}

func (l *Lexer) ScanText(input []rune) {
	l.scanMode = 't' // text file
	l.input = input
	l.advance() // Avanza al primer caracter
}

func (l *Lexer) ScanFile(fileName string) {
	l.scanMode = 'f' // file
	l.fileName = fileName
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("builder error: %s does not exist.", fileName)
	}
	// Convertimos el array de bytes a string
	// y luego a bytes runes
	l.input = []rune(string(fileContent))
	l.advance() // Avanza al primer caracter
}

func (l *Lexer) advance() {
	if l.ch == '\n' {
		l.line += 1
		l.col = 0
	} else {
		l.col += 1
	}
	if l.peekPos >= len(l.input) {
		l.ch = rune(0)
	} else {
		l.ch = l.input[l.peekPos]
	}
	l.pos = l.peekPos
	l.peekPos += 1
}

func (l *Lexer) peek() rune {
	if l.peekPos >= len(l.input) {
		return rune(0)
	}
	return l.input[l.peekPos]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // Ignora todos los espacios en blanco

	switch l.ch {
	case '\n':
		l.advance()
		tok = l.newToken(token.NewLine, "")
	case rune(0):
		if l.prevToken.Type != token.NewLine {
			tok = l.newToken(token.NewLine, "")
		} else {
			tok = l.newToken(token.Eof, "")
		}
	default:
		if l.isSymbol(l.ch) {
			peek := l.peek()
			key := string(l.ch) + string(peek)
			if t, ok := l.symbol[key]; ok { // Es un símbolo doble?
				l.advance() // avanza el primer símbolo
				l.advance() // avanza el segundo
				tok = l.newToken(t, key)
			} else {
				key = string(l.ch)
				if t, ok := l.symbol[key]; ok { // Es un símbolo sencillo?
					l.advance() // avanza el símbolo
					tok = l.newToken(t, key)
				} else {
					l.printError("invalid character literal.")
				}
			}
		} else if isDigit(l.ch) {
			lit := l.readNumber()
			tok = l.newToken(token.Number, lit)
		} else {
			l.advance()
			l.printError("invalid character literal.")
		}
	}
	return tok
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.advance()
	}
	return string(l.input[pos:l.pos])
}

func (l *Lexer) skipWhitespace() {
	for isSpace(l.ch) {
		if l.ch == '\n' && l.prevToken.Type != token.NewLine {
			break
		}
		l.advance()
	}
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\r' || ch == '\n' || ch == '\t'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isIdent(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}

func (l *Lexer) isSymbol(ch rune) bool {
	return strings.Contains(l.symbols, string(ch))
}

func (l *Lexer) printError(msg string) {
	if l.scanMode == 'f' {
		msg = fmt.Sprintf("%s:%d:%d: error: %s\n", l.fileName, l.line, l.col, msg)
	}
	fmt.Println(msg)
	os.Exit(1)
}
