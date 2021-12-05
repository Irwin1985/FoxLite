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
	prevToken token.TokenType
	symbol    map[string]token.TokenType
	symbols   string
}

func New() *Lexer {
	l := &Lexer{
		symbol:   map[string]token.TokenType{},
		symbols:  "+-*/^%=()[],¿?!<>.^",
		line:     1,
		col:      0,
		fileName: "",
	}
	l.prevToken = token.NewLine

	// Rellenamos los símbolos
	l.symbol["+"] = token.Plus
	l.symbol["+="] = token.PlusEq
	l.symbol["-"] = token.Minus
	l.symbol["-="] = token.MinusEq
	l.symbol["*"] = token.Mul
	l.symbol["*="] = token.MulEq
	l.symbol["/"] = token.Div
	l.symbol["/="] = token.DivEq
	l.symbol["="] = token.Assign
	l.symbol["=="] = token.Equal
	l.symbol["!="] = token.NotEq
	l.symbol["^"] = token.Pow
	l.symbol["%"] = token.Mod
	l.symbol["("] = token.Lparen
	l.symbol[")"] = token.Rparen
	l.symbol[","] = token.Comma
	l.symbol["["] = token.Lbracket
	l.symbol["]"] = token.Rbracket
	l.symbol["¿"] = token.OpenQM
	l.symbol["?"] = token.CloseQM
	l.symbol["!"] = token.Not
	l.symbol["<"] = token.Less
	l.symbol["<="] = token.LessEq
	l.symbol[">"] = token.Greater
	l.symbol[">="] = token.GreaterEq
	l.symbol["."] = token.Dot

	return l
}

func (l *Lexer) newToken(ttype token.TokenType, lit string, col int) token.Token {
	var t = token.Token{
		Type:    ttype,
		Literal: lit,
		Col:     col,
		Line:    l.line,
	}
	l.prevToken = ttype
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
	for !l.isAtEnd() {
		// ignorar espacios en blanco
		if isSpace(l.ch) {
			l.skipWhitespace()
			continue
		} // isSpace(l.ch)

		// ignorar comentarios
		if l.isComment() {
			l.skipComments()
			continue
		} // l.isComment()

		// identificadores
		if isLetter(l.ch) {
			col := l.col
			lit := l.readIdent()
			return l.newToken(token.LookupIdent(lit), lit, col)
		} // isIdent(l.ch)

		// números
		if isDigit(l.ch) {
			col := l.col
			lit := l.readNumber()
			return l.newToken(token.Number, lit, col)
		} // isDigit(l.ch)

		// string
		if isString(l.ch) {
			col := l.col
			return l.newToken(token.String, l.readString(), col)
		} // isString(l.ch)

		// salto de línea
		if l.ch == '\n' {
			col := l.col
			l.advance()
			return l.newToken(token.NewLine, "", col)
		} // l.ch == '\n'

		// caracteres especiales
		if l.isSymbol(l.ch) {
			peek := l.peek()
			key := string(l.ch) + string(peek)
			if t, ok := l.symbol[key]; ok { // Es un símbolo doble?
				col := l.col
				l.advance() // avanza el primer símbolo
				l.advance() // avanza el segundo
				return l.newToken(t, key, col)
			} else {
				key = string(l.ch)
				if t, ok := l.symbol[key]; ok { // Es un símbolo sencillo?
					col := l.col
					l.advance() // avanza el símbolo
					return l.newToken(t, key, col)
				} else {
					l.printError(fmt.Sprintf("invalid character literal [%s]", string(l.ch)))
				} // t, ok := l.symbol[key]; ok
			} // t, ok := l.symbol[key]; ok
		} // l.isSymbol(l.ch)
		tok := l.newToken(token.Illegal, string(l.ch), l.col)
		l.advance()
		return tok
	} // for l.ch != rune(0)
	// es EOF
	if l.prevToken != token.NewLine {
		return l.newToken(token.NewLine, "", 0)
	} else {
		return l.newToken(token.Eof, "", 0)
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) isSymbol(ch rune) bool {
	return strings.Contains(l.symbols, string(ch))
}

func (l *Lexer) isAtEnd() bool {
	return l.ch == rune(0)
}

func (l *Lexer) printError(msg string) {
	if l.scanMode == 'f' {
		msg = fmt.Sprintf("%s:%d:%d: error: %s\n", l.fileName, l.line, l.col, msg)
	}
	fmt.Println(msg)
	os.Exit(1)
}

func (l *Lexer) GetFileName() string {
	return l.fileName
}

func (l *Lexer) GetErrorFormat(t *token.Token) string {
	lincol := ""
	if t != nil {
		lincol = fmt.Sprintf("%d:%d", t.Line, t.Col)
	}
	if l.scanMode == 'f' {
		return fmt.Sprintf("%s:%s: error:", l.fileName, lincol)
	}
	return fmt.Sprintf("[%s]: error:", lincol)
}
