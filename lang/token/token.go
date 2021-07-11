package token

import "fmt"

type TokenType byte

const (
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// Identifiers + literals
	IDENT  // add, foobar, x, y, ...
	NUMBER // 1234
	STRING // "foobar"

	// Operators
	ASSIGN
	BINDING
	PLUS
	PLUS_EQ
	MINUS
	MINUS_EQ
	BANG
	MUL
	MUL_EQ
	DIV
	DIV_EQ

	// Comparison operators
	LT
	LEQ
	EQ
	NEQ
	GT
	GEQ

	// Delimiters
	COMMA
	DOT
	COLON
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// keywords
	FUNCTION
	FUNC
	LPARAMETERS
	ENDFUNC
	LOCAL
	PRIVATE
	PUBLIC
	TRUE
	FALSE
	NULL
	AND
	OR
	IF
	ELSE
	ENDIF
	DOCASE
	CASE
	OTHERWISE
	ENDCASE
	DOWHILE
	ENDDO
	FOR
	EACH
	TO
	STEP
	EXIT
	LOOP
	NEXT
	RETURN
)

// displayable tokens
var TokenNames = []string{
	"ILLEGAL",
	"EOF",
	"NEWLINE",

	"INDENT", // add, foobar, x, y, ...
	"NUMBER", // 1234
	"STRING", // "foobar"

	// Operators
	"ASSIGN",
	"BINDING",
	"PLUS",
	"PLUS_EQ",
	"MINUS",
	"MINUS_EQ",
	"BANG",
	"MUL",
	"MUL_EQ",
	"DIV",
	"DIV_EQ",

	// Comparison operators
	"LT",
	"LEQ",
	"EQ",
	"NEQ",
	"GT",
	"GEQ",

	// Delimiters
	"COMMA",
	"DOT",
	"COLON",
	"SEMICOLON",
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
	"LBRACKET",
	"RBRACKET",

	// keywords
	"FUNCTION",
	"FUNC",
	"LPARAMETERS",
	"ENDFUNC",
	"LOCAL",
	"PRIVATE",
	"PUBLIC",
	"TRUE",
	"FALSE",
	"NULL",
	"AND",
	"OR",
	"IF",
	"ELSE",
	"ENDIF",
	"DOCASE",
	"CASE",
	"OTHERWISE",
	"ENDCASE",
	"DOWHILE",
	"ENDDO",
	"FOR",
	"EACH",
	"TO",
	"STEP",
	"EXIT",
	"LOOP",
	"NEXT",
	"RETURN",
}

// keywords
var keywords = map[string]TokenType{
	"function":    FUNCTION,
	"func":        FUNCTION,
	"lparameters": LPARAMETERS,
	"endfunc":     ENDFUNC,
	"local":       LOCAL,
	"private":     PRIVATE,
	"public":      PUBLIC,
	".t.":         TRUE,
	".f.":         FALSE,
	"null":        NULL,
	"and":         AND,
	"or":          OR,
	"if":          IF,
	"else":        ELSE,
	"endif":       ENDIF,
	"return":      RETURN,
}

// small tokens
var smallTokens = map[string]TokenType{
	"=":  ASSIGN,
	":=": BINDING,
	"+":  PLUS,
	"+=": PLUS_EQ,
	"-":  MINUS,
	"-=": MINUS_EQ,
	"!":  BANG,
	"*":  MUL,
	"*=": MUL_EQ,
	"/":  DIV,
	"/=": DIV_EQ,

	// Comparison operators
	"<":  LT,
	"<=": LEQ,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GEQ,

	// Delimiters
	",": COMMA,
	":": COLON,
	";": SEMICOLON,
	"(": LPAREN,
	")": RPAREN,
	"{": LBRACE,
	"}": RBRACE,
	"[": LBRACKET,
	"]": RBRACKET,
}

// Token struct
type Token struct {
	Type   TokenType
	Lexeme interface{}
	Ln     int
	Col    int
}

func (t Token) ToString() string {
	return fmt.Sprintf("<Ln %d,\tCol %d\t%s,\t'%v'>", t.Ln, t.Col, TokenNames[t.Type], t.Lexeme)
}

// Return the matching token keyword or IDENT
func LookupIdent(key string) TokenType {
	if value, ok := keywords[key]; ok {
		return value
	}
	return IDENT
}

// Return the matching single or double small token
func Special(name string) (TokenType, bool) {
	if tok, ok := smallTokens[name]; ok {
		return tok, true
	}
	return EOF, false
}
