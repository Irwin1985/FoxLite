package lexer

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.advance()
	}
	return string(l.input[pos:l.pos])
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
