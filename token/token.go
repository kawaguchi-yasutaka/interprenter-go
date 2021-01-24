package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILEEGAL = "ILEEGAL"
	EOF     = "EOF"
)

const (
	IDENT = "IDENT"
	INT = "INT"
)

const (
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	SLASH = "/"
	ASTERISK = "*"

	LT = "<"
	RT = ">"
	BANG ="!"

	EQ = "=="
	NOT_EQ = "!="
)

const (
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
)

const (
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType{
	"let": LET,
	"fn": FUNCTION,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

func LookUpIdent(ident string) TokenType{
	if tok,ok := keywords[ident];ok{
		return tok
	}
	return IDENT
}
