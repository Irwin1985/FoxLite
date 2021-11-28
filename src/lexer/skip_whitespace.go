package lexer

func (l *Lexer) skipWhitespace() {
	for isSpace(l.ch) {
		l.advance()
	}
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\r' || ch == '\t'
}
