package main

import "fmt"

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR
	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER
	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

type TokenType = uint8

type Token struct {
	tokenType TokenType
	lexeme    []byte
	line      int
}

type Scanner struct {
	source  []byte
	start   int
	current int
	line    int
}

func NewScanner(source *[]byte) *Scanner {
	return &Scanner{
		source:  *source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanToken() Token {
	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF)
	}

	c := s.advance()

	switch c {
	case '(':
		return s.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return s.makeToken(TOKEN_SEMICOLON)
	case ',':
		return s.makeToken(TOKEN_COMMA)
	case '.':
		return s.makeToken(TOKEN_DOT)
	case '-':
		return s.makeToken(TOKEN_MINUS)
	case '+':
		return s.makeToken(TOKEN_PLUS)
	case '/':
		return s.makeToken(TOKEN_SLASH)
	case '*':
		return s.makeToken(TOKEN_STAR)
	case '!':
		return s.makeMatchedToken('=', TOKEN_BANG_EQUAL, TOKEN_BANG)
	case '=':
		return s.makeMatchedToken('=', TOKEN_EQUAL_EQUAL, TOKEN_EQUAL)
	case '<':
		return s.makeMatchedToken('=', TOKEN_LESS_EQUAL, TOKEN_LESS)
	case '>':
		return s.makeMatchedToken('=', TOKEN_GREATER_EQUAL, TOKEN_GREATER)
	}

	err := fmt.Sprintf("Unrecognized character, %v / \"%s\"", c, string(c))
	return s.errorToken(err)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    s.source[s.start:s.current],
		line:      s.line,
	}
}

func (s *Scanner) makeMatchedToken(expected byte, t TokenType, f TokenType) Token {
	if s.match(expected) {
		return s.makeToken(t)
	} else {
		return s.makeToken(f)
	}

}

func (s *Scanner) errorToken(message string) Token {
	return Token{
		tokenType: TOKEN_ERROR,
		lexeme:    []byte(message),
		line:      s.line,
	}
}
