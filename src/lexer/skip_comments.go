package lexer

func (l *Lexer) skipComments() {
	end := l.peek()
	if end == '/' { // avanzar hasta el final de la línea
		for !l.isAtEnd() && l.ch != '\n' {
			l.advance()
		}
		if l.ch == '\n' {
			l.advance()
		}
	} else if end == '*' { // avanzar hasta dar con la combinación '*/'
		for {
			if (l.ch == '*' && l.peek() == '/') || l.isAtEnd() {
				break
			}
			l.advance()
		}
		if l.isAtEnd() {
			l.printError("unterminated comment.")
		} else {
			l.advance() // avanza el '*'
			l.advance() // avanza el '/'
		}
	}
}

func (l *Lexer) isComment() bool {
	return l.ch == '/' && (l.peek() == '/' || l.peek() == '*')
}
