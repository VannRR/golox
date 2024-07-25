package lexer

import (
	"bytes"
	"golox/internal/token"
	"testing"
)

func Test_ScanToken(t *testing.T) {
	source := []byte(`( ) { } ; , . - + / * % ! = < > "hello" 123`)
	l := NewLexer(&source)

	expectedTokens := []token.TokenType{
		token.LeftParen, token.RightParen, token.LeftBrace, token.RightBrace,
		token.Semicolon, token.Comma, token.Dot, token.Minus, token.Plus,
		token.Slash, token.Star, token.Percent, token.Bang, token.Equal,
		token.Less, token.Greater, token.String, token.Number,
	}

	for _, expected := range expectedTokens {
		tok := l.ScanToken()
		if tok.Type != expected {
			t.Errorf("Expected token type %v, got %v", expected, tok.Type)
		}
	}
}

func Test_isAtEnd(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	if l.isAtEnd() {
		t.Error("Expected isAtEnd() to return false, but it returned true")
	}

	for !l.isAtEnd() {
		l.advance()
	}

	if !l.isAtEnd() {
		t.Error("Expected isAtEnd() to return true, but it returned false")
	}
}

func Test_advance(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	expectedByte := byte('s')
	actualByte := l.advance()
	if actualByte != expectedByte {
		t.Errorf("Expected byte %c, but got %c", expectedByte, actualByte)
	}
}

func Test_peek(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	expectedByte := byte('s')
	actualByte := l.peek()
	if actualByte != expectedByte {
		t.Errorf("Expected peeked byte %c, but got %c", expectedByte, actualByte)
	}
}

func Test_peekNext(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	expectedByte := byte('o')
	actualByte := l.peekNext()
	if actualByte != expectedByte {
		t.Errorf("Expected peeked next byte %c, but got %c", expectedByte, actualByte)
	}
}

func Test_match(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	if !l.match('s') {
		t.Error("Expected match('s') to return true, but it returned false")
	}

	if l.match('x') {
		t.Error("Expected match('x') to return false, but it returned true")
	}
}

func Test_makeToken(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	expectedType := token.Identifier
	tok := l.makeToken(expectedType)
	if tok.Type != expectedType {
		t.Errorf("Expected token type %v, but got %v", expectedType, tok.Type)
	}
}

func Test_makeMatchedToken(t *testing.T) {
	source := []byte("some source code")
	l := NewLexer(&source)

	expectedType := token.Equal
	tok := l.makeMatchedToken('=', token.EqualEqual, expectedType)
	if tok.Type != expectedType {
		t.Errorf("Expected token type %v, but got %v", expectedType, tok.Type)
	}
}

func Test_skipWhitespace(t *testing.T) {
	source := []byte("  \t\n\r some source code")
	l := NewLexer(&source)

	tok := l.skipWhitespace()

	if tok.Type != 0 {
		t.Errorf("Expected token type Identifier, but got %v", tok.Type)
	}
}

func Test_skipBlockComment(t *testing.T) {
	source := []byte("/* This is a block comment */ some source code")
	l := NewLexer(&source)

	l.advance()

	tok := l.skipBlockComment()

	if tok.Type != 0 {
		t.Errorf("Expected token type Error, but got %v", tok.Type)
	}
}

func Test_string(t *testing.T) {
	source := []byte(`"Hello, world!" some source code`)
	l := NewLexer(&source)

	l.advance()

	tok := l.string()

	if tok.Type != token.String {
		t.Errorf("Expected token type String, but got %v", tok.Type)
	}

	expectedLexeme := []byte(`"Hello, world!"`)
	if !bytes.Equal(tok.Lexeme, expectedLexeme) {
		t.Errorf("Expected lexeme %q, but got %q", expectedLexeme, tok.Lexeme)
	}
}

func Test_number(t *testing.T) {
	source := []byte("123.456")
	l := NewLexer(&source)

	l.advance()

	tok := l.number()

	if tok.Type != token.Number {
		t.Errorf("Expected token type Number, but got %v", tok.Type)
	}
}

func Test_identifierType(t *testing.T) {
	source := []byte("if and else while;")
	l := NewLexer(&source)

	expectedTypes := []token.TokenType{
		token.If, token.And, token.Else, token.While,
	}

	for _, expected := range expectedTypes {
		tokType := l.ScanToken().Type
		if tokType != expected {
			t.Errorf("Expected token type %v, but got %v", expected, tokType)
		}
	}
}

func Test_checkKeyword(t *testing.T) {
	source := []byte("var class true;")
	l := NewLexer(&source)

	keywords := [][]byte{
		[]byte("var"), []byte("class"), []byte("true"),
	}

	for _, expected := range keywords {
		keyword := l.ScanToken().Lexeme
		if string(keyword) != string(expected) {
			t.Errorf("Expected keyword %v, but got %v", string(expected), string(keyword))
		}
	}
}

func Test_isAlpha(t *testing.T) {
	if !isAlpha('a') {
		t.Error("Expected 'a' to be alpha, but it's not.")
	}

	if !isAlpha('Z') {
		t.Error("Expected 'Z' to be alpha, but it's not.")
	}

	if !isAlpha('_') {
		t.Error("Expected '_' to be alpha, but it's not.")
	}

	nonAlphaChars := []byte{'1', ' ', '*', '@', '$'}
	for _, c := range nonAlphaChars {
		if isAlpha(c) {
			t.Errorf("Expected '%c' not to be alpha, but it is.", c)
		}
	}
}

func Test_isDigit(t *testing.T) {
	for d := '0'; d <= '9'; d++ {
		if !isDigit(byte(d)) {
			t.Errorf("Expected '%c' to be a digit, but it's not.", d)
		}
	}

	nonDigitChars := []byte{'a', 'Z', '_', ' ', '*', '@', '$'}
	for _, c := range nonDigitChars {
		if isDigit(c) {
			t.Errorf("Expected '%c' not to be a digit, but it is.", c)
		}
	}
}
