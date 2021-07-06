package token

import "fmt"

type TokenType byte

const (
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// Identifiers + literals
	IDENT  // add, foobar, x, y, ...
	INT    // 1234
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
var tokenNames = []string{
	"ILLEGAL",
	"EOF",
	"NEWLINE",

	"INDENT", // add, foobar, x, y, ...
	"INT",    // 1234
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
	"function": FUNCTION,
	"endfunc":  ENDFUNC,
	"local":    LOCAL,
	"private":  PRIVATE,
	"public":   PUBLIC,
	".t.":      TRUE,
	".f.":      FALSE,
	"null":     NULL,
	"and":      AND,
	"or":       OR,
	"if":       IF,
	"else":     ELSE,
	"endif":    ENDIF,
	"return":   RETURN,
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

// boolprecursor
// '=', '==', '<', '<=', '>', '>=', 'and', 'or', '!', '!='
var boolprecursor = map[TokenType]byte{
	ASSIGN:  0,
	EQ:      1,
	LT:      2,
	LEQ:     3,
	GT:      4,
	GEQ:     5,
	AND:     6,
	OR:      7,
	BANG:    8,
	NEQ:     9,
	EOF:     10,
	NEWLINE: 11,
}

// Token struct
type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Col    int
}

func (t *Token) ToString() string {
	return fmt.Sprintf("<Ln %d, Col %d \t <%s, '%s'>", t.Line, t.Col, tokenNames[t.Type], t.Lexeme)
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

// the given token is a boolean precursor?
func BoolPrecursor(t TokenType) bool {
	if _, ok := boolprecursor[t]; ok {
		return true
	}
	return false
}
