package lox

import "fmt"

// TokenType includes every type of operator, keyword, and literal in Lox
type TokenType uint8

// TokenType values
const (
	// single character
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// 1-2 characters
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

// Token represents a token as produced by the lexer. Lexeme stores the string value of the token, and line the line number of the original file where the token is located
type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

// NewToken creates a new token of the given type
func NewToken(typ TokenType, lexeme string, line int) *Token {
	return &Token{typ, lexeme, line}
}

// ToString represents a token as a string
func (t Token) ToString() string {
	return fmt.Sprintf("%v %v", t.Type, t.Lexeme)
}
