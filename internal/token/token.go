package token

import "fmt"

const (
	// Single-character tokens.
	LeftParen Type = iota
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
	Percent
	// One or two character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	// Literals.
	Identifier
	String
	Number
	// Keywords.
	And
	Class
	Else
	False
	For
	Fun
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

	Error
	Eof
)

type Type = uint8

type Token struct {
	Type   Type
	Lexeme []byte
	Line   uint16
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %v, Lexeme: %s, Line: %d", t.Type, t.Lexeme, t.Line)
}
