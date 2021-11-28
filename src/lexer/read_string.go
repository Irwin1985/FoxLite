package lexer

func (l *Lexer) readString() string {
	end := l.ch
	pos := l.pos + 1
	for {
		l.advance()
		if l.ch == end || l.isAtEnd() {
			break
		}
	}
	if l.isAtEnd() {
		l.printError("unfinished string literal")
	}
	str := string(l.input[pos:l.pos])
	l.advance() // avanza el cierre del string
	return str
}

func isString(ch rune) bool {
	return ch == '"' || ch == '\'' || ch == '`'
}
