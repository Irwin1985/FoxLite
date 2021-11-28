package lexer

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isIdent(l.ch) {
		l.advance()
	}
	return string(l.input[pos:l.pos])
}

func isIdent(ch rune) bool {
	return isLetter(ch) || isDigit(ch)
}
