package compiler

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/lexer"
	"golox/internal/opcode"
	"golox/internal/token"
	"golox/internal/value"
	"testing"
)

func Test_GetRule(t *testing.T) {
	testCases := []struct {
		tkn      token.TokenType
		expected *ParseRule
	}{
		{
			tkn:      token.Plus,
			expected: &ParseRule{prefix: nil, infix: (*Parser).binary, precedence: PrecTerm},
		},
		{
			tkn:      token.LeftParen,
			expected: &ParseRule{prefix: (*Parser).grouping, infix: nil, precedence: PrecNone},
		}, {
			tkn:      token.Slash,
			expected: &ParseRule{prefix: nil, infix: (*Parser).binary, precedence: PrecFactor},
		},
	}

	for _, tc := range testCases {
		r := getRule(tc.tkn)
		if fmt.Sprint(r) != fmt.Sprint(tc.expected) {
			t.Errorf("Expected ParseRule '%v', got '%v'.", tc.expected, r)
		}
	}
}

func Test_NewParser(t *testing.T) {
	s := []byte("var foo = 1;")
	l := lexer.NewLexer(&s)
	c := chunk.NewChunk()

	p := NewParser(l, c)

	expected := &Parser{
		lexer:     l,
		chunk:     c,
		current:   token.Token{},
		previous:  token.Token{},
		hadError:  false,
		panicMode: false,
	}

	if fmt.Sprint(p) != fmt.Sprint(expected) {
		t.Errorf("Expected NewParser '%v', got '%v'.", expected, p)
	}
}

func Test_Compile(t *testing.T) {
	s := []byte("var foo = (1 / 0.3) + (20 - 2) * 11; var bar = foo % 3;")
	c := chunk.NewChunk()

	ok := Compile(&s, c)

	expectedCode := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.Constant, 2,
		opcode.Divide,
		opcode.Constant, 3,
		opcode.Constant, 4,
		opcode.Subtract,
		opcode.Constant, 5,
		opcode.Multiply,
		opcode.Add,
		opcode.DefineGlobal, 0,
		opcode.Constant, 6,
		opcode.Constant, 7,
		opcode.GetGlobal, 7,
		opcode.Constant, 8,
		opcode.Modulo,
		opcode.DefineGlobal, 6,
		opcode.Return,
	}

	expectedConstants := []value.Value{
		value.StringVal("foo"),
		value.NumberVal(1),
		value.NumberVal(0.3),
		value.NumberVal(20),
		value.NumberVal(2),
		value.NumberVal(11),
		value.StringVal("bar"),
		value.StringVal("foo"),
		value.NumberVal(3),
	}

	if ok != true {
		t.Errorf("Expected Compile to return 'true' to indicate no errors.")
	}

	checkOpcodes(t, c.Code, expectedCode)

	checkConstants(t, c.Constants, expectedConstants)
}

func Test_printStatement(t *testing.T) {
	input := "print"
	p := setupParserForTest(input)

	p.printStatement()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.GetGlobal, 0,
		opcode.Print,
	}

	expectedConstants := []value.Value{
		value.StringVal(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_varDeclaration(t *testing.T) {
	p := setupParserForTest("var foo;")

	p.varDeclaration()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Nil,
		opcode.DefineGlobal, 0,
	}

	expectedConstants := []value.Value{
		value.StringVal(""),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_expression(t *testing.T) {
	p := setupParserForTest("2 -3")

	p.expression()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.Subtract,
	}

	expectedConstants := []value.Value{
		value.NumberVal(2),
		value.NumberVal(3),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_grouping(t *testing.T) {
	p := setupParserForTest("(1 + 2)")

	p.grouping()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.Add,
	}

	expectedConstants := []value.Value{
		value.NumberVal(1),
		value.NumberVal(2),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_number(t *testing.T) {
	input := 420
	p := setupParserForTest("")
	p.previous.Lexeme = []byte(fmt.Sprint(input))

	p.number()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
	}

	expectedConstants := []value.Value{
		value.NumberVal(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_variable(t *testing.T) {
	input := "wow"
	s := []byte(input)
	l := lexer.NewLexer(&s)
	c := chunk.NewChunk()
	p := NewParser(l, c)
	p.previous = l.ScanToken()

	p.variable()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.GetGlobal, 0,
	}

	expectedConstants := []value.Value{
		value.StringVal(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_string(t *testing.T) {
	p := setupParserForTest("")
	p.previous.Lexeme = []byte("\"wow\"")

	p.string()

	expectedOpcodes := []byte{
		opcode.Constant, 0,
	}

	expectedConstants := []value.Value{
		value.StringVal("wow"),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_binary(t *testing.T) {
	table := []struct {
		t token.TokenType
		b byte
	}{
		{token.BangEqual, opcode.NotEqual},
		{token.EqualEqual, opcode.Equal},
		{token.Greater, opcode.Greater},
		{token.GreaterEqual, opcode.GreaterEqual},
		{token.Less, opcode.Less},
		{token.LessEqual, opcode.LessEqual},
		{token.Plus, opcode.Add},
		{token.Minus, opcode.Subtract},
		{token.Star, opcode.Multiply},
		{token.Slash, opcode.Divide},
		{token.Percent, opcode.Modulo},
	}

	for _, pair := range table {

		p := setupParserForTest("")

		p.previous.Type = pair.t

		p.binary()

		expectedOpcodes := []byte{
			pair.b,
		}

		expectedConstants := []value.Value{}

		checkOpcodes(t, p.chunk.Code, expectedOpcodes)

		checkConstants(t, p.chunk.Constants, expectedConstants)
	}
}

func Test_literal(t *testing.T) {
	table := []struct {
		t token.TokenType
		b byte
	}{
		{token.False, opcode.False},
		{token.Nil, opcode.Nil},
		{token.True, opcode.True},
	}

	for _, pair := range table {

		p := setupParserForTest("")

		p.previous.Type = pair.t

		p.literal()

		expectedOpcodes := []byte{
			pair.b,
		}

		expectedConstants := []value.Value{}

		checkOpcodes(t, p.chunk.Code, expectedOpcodes)

		checkConstants(t, p.chunk.Constants, expectedConstants)
	}
}

func Test_unary(t *testing.T) {
	table := []struct {
		t token.TokenType
		b byte
	}{
		{token.Bang, opcode.Not},
		{token.Minus, opcode.Negate},
	}

	for _, pair := range table {

		p := setupParserForTest("")

		p.previous.Type = pair.t

		p.unary()

		expectedOpcodes := []byte{
			pair.b,
		}

		expectedConstants := []value.Value{}

		checkOpcodes(t, p.chunk.Code, expectedOpcodes)

		checkConstants(t, p.chunk.Constants, expectedConstants)
	}
}

func Test_parsePrecedence(t *testing.T) {
	p := setupParserForTest("1 + 2")

	p.parsePrecedence(PrecTerm)

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.Add,
	}

	expectedConstants := []value.Value{
		value.NumberVal(1),
		value.NumberVal(2),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_parseVariable(t *testing.T) {
	input := "wow"
	p := setupParserForTest(input)

	p.current = p.lexer.ScanToken()

	p.parseVariable([]byte("this is a test"))

	expectedOpcodes := []byte{
		opcode.Constant, 0,
	}

	expectedConstants := []value.Value{
		value.StringVal(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_defineVariable(t *testing.T) {
	index := 0
	p := setupParserForTest("wow")

	p.current = p.lexer.ScanToken()

	p.defineVariable(index)

	expectedOpcodes := []byte{
		opcode.DefineGlobal, byte(index),
	}

	expectedConstants := []value.Value{}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_error(t *testing.T) {
	p := setupParserForTest("wow")

	p.current = p.lexer.ScanToken()

	p.error([]byte("this is a test error msg"))

	if p.panicMode != true {
		t.Error("Expected Parser panicMode to be true")
	}

	if p.hadError != true {
		t.Error("Expected Parser hadError to be true")
	}
}

func Test_match(t *testing.T) {
	p := setupParserForTest("wow")

	p.current = p.lexer.ScanToken()

	p.error([]byte("this is a test error msg"))

	if p.panicMode != false {
		t.Error("TODO: add test for match method")
	}
}

func checkOpcodes(t *testing.T, got []byte, expected []byte) {
	t.Helper()

	gotLen := len(got)
	expectedLen := len(expected)

	if gotLen != expectedLen {
		t.Fatalf("Expected byte code slice with length '%v', got '%v'.", expectedLen, gotLen)
	}

	for i := 0; i < gotLen; i++ {
		switch expected[i] {
		case opcode.Constant, opcode.DefineGlobal, opcode.GetGlobal:
			if got[i] != expected[i] {
				t.Errorf("Expected opcode '%v' at index %v, got '%v'.", opcode.Name(expected[i]), i, opcode.Name(got[i]))
			}
			i++
			if got[i] != expected[i] {
				t.Errorf("Expected constant index '%v', got '%v'.", expected[i], got[i])
			}
		default:
			if got[i] != expected[i] {
				t.Errorf("Expected opcode '%v' at index %v, got '%v'.", opcode.Name(expected[i]), i, opcode.Name(got[i]))
			}
		}
	}
}

func setupParserForTest(source string) *Parser {
	s := []byte(source)
	l := lexer.NewLexer(&s)
	c := chunk.NewChunk()
	return NewParser(l, c)
}

func checkConstants(t *testing.T, got []value.Value, expected []value.Value) {
	t.Helper()

	gotLen := len(got)
	expectedLen := len(expected)

	if gotLen != expectedLen {
		t.Fatalf("Expected constants slice with length '%v', got '%v'.", expectedLen, gotLen)
	}

	for i := 0; i < gotLen; i++ {
		constant := got[i]
		expectedConstant := expected[i]
		if constant != expectedConstant {
			t.Errorf("Expected constant '%v' at index %v, got '%v'.", expectedConstant, i, constant)
		}
	}

}
