package lexer

const EOF_CHAR rune = rune(0)

type Stream struct {
	input []rune
	pos   int
	Line  int
	Col   int
}

func NewStream(input string) *Stream {
	s := &Stream{
		input: []rune(input),
		pos:   -1,
		Line:  1,
		Col:   0,
	}
	return s
}

func (s *Stream) Read() rune {
	var c rune
	s.pos += 1
	if s.pos >= len(s.input) {
		return EOF_CHAR
	}
	c = s.input[s.pos]
	if c == rune('\n') {
		s.Line += 1
		s.Col = 0
	} else {
		s.Col += 1
	}
	return c
}

func (s *Stream) Peek() rune {
	if s.End() {
		return EOF_CHAR
	}
	return s.input[s.pos+1]
}

func (s *Stream) End() bool {
	return s.pos+1 >= len(s.input)
}
