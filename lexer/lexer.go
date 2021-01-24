package lexer

import (
	"interpreter-go/token"
)

type Lexer struct {
	input string
	position int
	reandPosition int
	ch byte
}

func New(input string)Lexer {
	lexer := Lexer{
		input: input,
	}
	lexer.readChar()
	return lexer
}

func (l *Lexer)NextToken() token.Token{
	l.skipWhiteSpace()
	var tok token.Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Literal: literal,Type: token.EQ}
		}else {
			tok = newToken(token.ASSIGN,l.ch)
		}
	case '+':
		tok = newToken(token.PLUS,l.ch)
	case '(':
		tok = newToken(token.LPAREN,l.ch)
	case ')':
		tok = newToken(token.RPAREN,l.ch)
	case '{':
		tok = newToken(token.LBRACE,l.ch)
	case '}':
		tok = newToken(token.RBRACE,l.ch)
	case ',':
		tok = newToken(token.COMMA,l.ch)
	case ';':
		tok = newToken(token.SEMICOLON,l.ch)
	case '-':
		tok = newToken(token.MINUS,l.ch)
	case '/':
		tok = newToken(token.SLASH,l.ch)
	case '*':
		tok = newToken(token.ASTERISK,l.ch)
	case '<':
		tok = newToken(token.LT,l.ch)
	case '>':
		tok = newToken(token.RT,l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Literal: literal,Type: token.NOT_EQ}
		}else {
			tok = newToken(token.BANG,l.ch)
		}
	case 0:
		tok = newToken(token.EOF,l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch){
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILEEGAL,l.ch)
		}
	}
	l.readChar()
	return tok
}


func (l *Lexer) readChar() {
	if l.reandPosition >= len(l.input) {
		l.ch = 0
	}else {
		l.ch = l.input[l.reandPosition]
	}
	l.position = l.reandPosition
	l.reandPosition = l.reandPosition + 1
}


func (l *Lexer)readIdentifier() string{
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer)skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r'{
		l.readChar()
	}
}

func (l *Lexer)readNumber() string{
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l Lexer)peekChar() byte{
	if l.reandPosition >= len(l.input) {
		return '0'
	}else {
		return l.input[l.reandPosition]
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool{
	return 'a' <= ch && 'z' >= ch || 'A' <= ch && 'Z' >= ch || '_' == ch
}

func newToken(tokenType token.TokenType,ch byte) token.Token {
	return token.Token{Type: tokenType,Literal: string(ch)}
}


