package lexer

import (
	"bytes"
	"fmt"
	"golox/internal/token"
)

type Lexer struct {
	source  []byte
	start   int
	current int
	line    uint16
}

func NewLexer(source *[]byte) *Lexer {
	return &Lexer{
		source:  *source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (l *Lexer) ScanToken() token.Token {
	t := l.skipWhitespace()
	if t.Type == token.Error {
		return t
	}
	l.start = l.current

	if l.isAtEnd() {
		return l.makeToken(token.Eof)
	}

	c := l.advance()

	if isAlpha(c) {
		return l.identifier()
	}

	if isDigit(c) {
		return l.number()
	}

	switch c {
	case '(':
		return l.makeToken(token.LeftParen)
	case ')':
		return l.makeToken(token.RightParen)
	case '{':
		return l.makeToken(token.LeftBrace)
	case '}':
		return l.makeToken(token.RightBrace)
	case ';':
		return l.makeToken(token.Semicolon)
	case ',':
		return l.makeToken(token.Comma)
	case '.':
		return l.makeToken(token.Dot)
	case '-':
		return l.makeToken(token.Minus)
	case '+':
		return l.makeToken(token.Plus)
	case '/':
		return l.makeToken(token.Slash)
	case '*':
		return l.makeToken(token.Star)
	case '%':
		return l.makeToken(token.Percent)
	case '!':
		return l.makeMatchedToken('=', token.BangEqual, token.Bang)
	case '=':
		return l.makeMatchedToken('=', token.EqualEqual, token.Equal)
	case '<':
		return l.makeMatchedToken('=', token.LessEqual, token.Less)
	case '>':
		return l.makeMatchedToken('=', token.GreaterEqual, token.Greater)
	case '"':
		return l.string()
	}

	err := fmt.Sprintf("Unrecognized character, %v / \"%s\"", c, string(c))
	return l.errorToken(err)
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) makeToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type:   tokenType,
		Lexeme: l.source[l.start:l.current],
		Line:   l.line,
	}
}

func (l *Lexer) makeMatchedToken(expected byte, t token.TokenType, f token.TokenType) token.Token {
	if l.match(expected) {
		return l.makeToken(t)
	} else {
		return l.makeToken(f)
	}

}

func (l *Lexer) errorToken(message string) token.Token {
	return token.Token{
		Type:   token.Error,
		Lexeme: []byte(message),
		Line:   l.line,
	}
}

func (l *Lexer) skipWhitespace() token.Token {
	for {
		switch c := l.peek(); c {
		case ' ', '\r', '\t':
			l.current++
		case '\n':
			l.line++
			l.current++
		case '/':
			switch nc := l.peekNext(); nc {
			case '/':
				for l.peek() != '\n' && !l.isAtEnd() {
					l.current++
				}
			case '*':
				l.current += 2
				t := l.skipBlockComment()
				if t.Type == token.Error {
					return t
				}
				return l.skipWhitespace()
			default:
				return token.Token{}
			}
		default:
			return token.Token{}
		}
	}
}

func (l *Lexer) skipBlockComment() token.Token {
	for !l.isAtEnd() {
		if l.peek() == '*' && l.peekNext() == '/' {
			l.current += 2
			return token.Token{}
		} else if l.peek() == '/' && l.peekNext() == '*' {
			l.current += 2
			l.skipBlockComment()
		} else {
			l.current++
		}
	}

	if l.isAtEnd() {
		return l.errorToken("Unterminated block comment.")
	}

	return token.Token{}
}

func (l *Lexer) string() token.Token {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.current++
	}

	if l.isAtEnd() {
		return l.errorToken("Unterminated string.")
	}

	l.current++
	return l.makeToken(token.String)
}

func (l *Lexer) identifier() token.Token {
	for isAlpha(l.peek()) || isDigit(l.peek()) {
		l.current++
	}
	return l.makeToken(l.identifierType())
}

func (l *Lexer) number() token.Token {
	for isDigit(l.peek()) {
		l.current++
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.current++
		for isDigit(l.peek()) {
			l.current++
		}
	}

	return l.makeToken(token.Number)
}

func (l *Lexer) identifierType() token.TokenType {
	switch c := l.source[l.start]; c {
	case 'a':
		return l.checkKeyword(1, []byte("nd"), token.And)
	case 'c':
		return l.checkKeyword(1, []byte("lass"), token.Class)
	case 'e':
		return l.checkKeyword(1, []byte("lse"), token.Else)
	case 'f':
		if l.current-l.start > 1 {
			switch nc := l.source[l.start+1]; nc {
			case 'a':
				return l.checkKeyword(2, []byte("lse"), token.False)
			case 'o':
				return l.checkKeyword(2, []byte("r"), token.For)
			case 'u':
				return l.checkKeyword(2, []byte("n"), token.Fun)
			}
		}
	case 'i':
		return l.checkKeyword(1, []byte("f"), token.If)
	case 'n':
		return l.checkKeyword(1, []byte("il"), token.Nil)
	case 'o':
		return l.checkKeyword(1, []byte("r"), token.Or)
	case 'p':
		return l.checkKeyword(1, []byte("rint"), token.Print)
	case 'r':
		return l.checkKeyword(1, []byte("eturn"), token.Return)
	case 's':
		return l.checkKeyword(1, []byte("uper"), token.Super)
	case 't':
		if l.current-l.start > 1 {
			switch nc := l.source[l.start+1]; nc {
			case 'h':
				return l.checkKeyword(2, []byte("is"), token.This)
			case 'r':
				return l.checkKeyword(2, []byte("ue"), token.True)
			}
		}
	case 'v':
		return l.checkKeyword(1, []byte("ar"), token.Var)
	case 'w':
		return l.checkKeyword(1, []byte("hile"), token.While)
	}

	return token.Identifier
}

func (l *Lexer) checkKeyword(start int, rest []byte, t token.TokenType) token.TokenType {
	s := l.start + start
	e := len(rest) + s
	if e < len(l.source) && bytes.Equal(l.source[s:e], rest) {
		return t
	}
	return token.Identifier
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
