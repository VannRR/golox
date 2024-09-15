package compiler

import (
	"fmt"
	"github.com/VannRR/golox/internal/chunk"
	"github.com/VannRR/golox/internal/lexer"
	"github.com/VannRR/golox/internal/object"
	"github.com/VannRR/golox/internal/opcode"
	"github.com/VannRR/golox/internal/token"
	"github.com/VannRR/golox/internal/value"
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

func Test_NewCompiler(t *testing.T) {
	co := NewCompiler()

	expected := &Compiler{
		locals:     make([]Local, 1),
		localCount: 0,
		scopeDepth: 0,
	}

	if fmt.Sprint(co) != fmt.Sprint(expected) {
		t.Errorf("Expected NewCompiler '%v', got '%v'.", expected, co)
	}
}

func Test_NewParser(t *testing.T) {
	s := []byte("var foo = 1;")
	l := lexer.NewLexer(&s)
	co := NewCompiler()
	ch := chunk.NewChunk()

	p := NewParser(l, ch, co)

	expected := &Parser{
		lexer:     l,
		chunk:     ch,
		compiler:  co,
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
		object.ObjString("foo"),
		value.NumberVal(1),
		value.NumberVal(0.3),
		value.NumberVal(20),
		value.NumberVal(2),
		value.NumberVal(11),
		object.ObjString("bar"),
		object.ObjString("foo"),
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
		object.ObjString(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_forStatement(t *testing.T) {
	p := setupParserForTest("")

	p.forStatement()

	expectedOpcodes := []byte{
		opcode.Pop,
		opcode.JumpIfFalse, 0, 12,
		opcode.Pop,
		opcode.Jump, 0, 4,
		opcode.Pop,
		opcode.Loop, 0, 11,
		opcode.Pop,
		opcode.Loop, 0, 8,
		opcode.Pop,
	}

	expectedConstants := []value.Value{}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_ifStatement(t *testing.T) {
	p := setupParserForTest("")

	p.ifStatement()

	expectedOpcodes := []byte{
		opcode.JumpIfFalse, 0, 5,
		opcode.Pop,
		opcode.Pop,
		opcode.Jump, 0, 1,
		opcode.Pop,
	}

	expectedConstants := []value.Value{}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_whileStatement(t *testing.T) {
	p := setupParserForTest("while (false) print 1;")

	p.whileStatement()

	expectedOpcodes := []byte{
		opcode.JumpIfFalse, 0, 6,
		opcode.Pop,
		opcode.False,
		opcode.Pop,
		opcode.Loop, 0, 9,
		opcode.Pop,
	}

	expectedConstants := []value.Value{}

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
		object.ObjString(""),
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

func Test_block(t *testing.T) {
	p := setupParserForTest("{var foo = 1;}")

	p.block()

	expectedOpcodes := []byte{
		opcode.Pop,
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.DefineGlobal, 0,
	}

	expectedConstants := []value.Value{
		object.ObjString("foo"),
		value.NumberVal(1),
	}

	if p.panicMode {
		t.Error("Expected no panic from block")
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_block_fail(t *testing.T) {
	p := setupParserForTest("{var foo = 1;")

	p.block()

	if !p.panicMode {
		t.Error("Expected panic from block without '}'.")
	}
}

func Test_grouping(t *testing.T) {
	p := setupParserForTest("(1 + 2)")

	p.grouping(false)

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

	p.number(false)

	expectedOpcodes := []byte{
		opcode.Constant, 0,
	}

	expectedConstants := []value.Value{
		value.NumberVal(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_or(t *testing.T) {
	p := setupParserForTest("")

	p.or(false)

	expectedOpcodes := []byte{
		opcode.JumpIfFalse, 0, 3,
		opcode.Jump, 0, 1,
		opcode.Pop,
	}

	expectedConstants := []value.Value{}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_variable_get(t *testing.T) {
	s := []byte(`var wow = 1; var foo = wow + 1;`)
	c := chunk.NewChunk()
	Compile(&s, c)

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.DefineGlobal, 0,
		opcode.Constant, 2,
		opcode.Constant, 3,
		opcode.GetGlobal, 3,
		opcode.Constant, 4,
		opcode.Add,
		opcode.DefineGlobal, 2,
		opcode.Return,
	}

	expectedConstants := []value.Value{
		object.ObjString("wow"),
		value.NumberVal(1),
		object.ObjString("foo"),
		object.ObjString("wow"),
		value.NumberVal(1),
	}

	checkOpcodes(t, c.Code, expectedOpcodes)

	checkConstants(t, c.Constants, expectedConstants)
}

func Test_variable_set(t *testing.T) {
	s := []byte(`var wow = 1; wow = 2;`)
	c := chunk.NewChunk()
	Compile(&s, c)

	expectedOpcodes := []byte{
		opcode.Constant, 0,
		opcode.Constant, 1,
		opcode.DefineGlobal, 0,
		opcode.Constant, 2,
		opcode.Constant, 3,
		opcode.SetGlobal, 2,
		opcode.Pop,
		opcode.Return,
	}

	expectedConstants := []value.Value{
		object.ObjString("wow"),
		value.NumberVal(1),
		object.ObjString("wow"),
		value.NumberVal(2),
	}

	checkOpcodes(t, c.Code, expectedOpcodes)

	checkConstants(t, c.Constants, expectedConstants)
}

func Test_namedVariable(t *testing.T) {
	p := setupParserForTest("")

	p.namedVariable(token.Token{Type: token.Identifier, Lexeme: []byte("myVar")}, true)
	expectedOne := []byte{
		opcode.Constant, 0,
		opcode.GetGlobal, 0,
	}
	checkOpcodes(t, p.chunk.Code, expectedOne)

	p.namedVariable(token.Token{Type: token.Identifier, Lexeme: []byte("anotherVar")}, false)
	expectedTwo := []byte{
		opcode.Constant, 0,
		opcode.GetGlobal, 0,
		opcode.Constant, 1,
		opcode.GetGlobal, 1,
	}
	checkOpcodes(t, p.chunk.Code, expectedTwo)
}

func Test_resolveLocal(t *testing.T) {
	p := setupParserForTest("")

	l := []byte("mylocal")

	p.compiler.localCount = 2

	p.compiler.locals = append(
		p.compiler.locals,
		Local{depth: 0, name: token.Token{Type: token.Identifier, Lexeme: l}},
	)

	index := p.resolveLocal(&token.Token{Type: token.Identifier, Lexeme: l})
	if index != 1 {
		t.Errorf("Expected index == 1, got %v", index)
	}

	index = p.resolveLocal(&token.Token{Type: token.Identifier, Lexeme: []byte("nonExistent")})
	if index != -1 {
		t.Errorf("Expected index == -1, got %v", index)
	}
}

func Test_string(t *testing.T) {
	p := setupParserForTest("")
	p.previous.Lexeme = []byte("\"wow\"")

	p.string(false)

	expectedOpcodes := []byte{
		opcode.Constant, 0,
	}

	expectedConstants := []value.Value{
		object.ObjString("wow"),
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

		p.binary(false)

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

		p.literal(false)

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

		p.unary(false)

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
		object.ObjString(input),
	}

	checkOpcodes(t, p.chunk.Code, expectedOpcodes)

	checkConstants(t, p.chunk.Constants, expectedConstants)
}

func Test_parseVariable_ScopeDepthGreaterThanZero(t *testing.T) {
	input := "anotherVar"
	p := setupParserForTest(input)

	p.current = p.lexer.ScanToken()
	p.compiler.scopeDepth = 1

	result := p.parseVariable([]byte("Scope depth test"))
	if result != 0 {
		t.Errorf("Expected result from parseVariable to be 0, got %v", result)
	}
}

func Test_addLocal(t *testing.T) {
	p := setupParserForTest("")

	localVarName := token.Token{Type: token.Identifier, Lexeme: []byte("myVar")}
	p.compiler.localCount = 0
	p.addLocal(localVarName)

	if p.compiler.localCount != 1 {
		t.Errorf("Expected compiler localCount to be 1, got %v", p.compiler.localCount)
	}

	p.compiler.localCount = 999999999
	p.addLocal(token.Token{Type: token.Identifier, Lexeme: []byte("tooManyVar")})

	if p.hadError != true {
		t.Error("Expected error from trying to add to many local variables.")
	}
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

func Test_and(t *testing.T) {
	p := setupParserForTest("")

	p.and(false)

	expectedOpcodes := []byte{
		opcode.JumpIfFalse, 0, 1,
		opcode.Pop,
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
	p := setupParserForTest("123")

	numberToken := p.lexer.ScanToken()

	p.current = numberToken

	if p.match(token.Number) != true {
		t.Error("Expected token to match type Number")
	}

	p.current = numberToken

	if p.match(token.Nil) != false {
		t.Error("Expected token to not match type Nil")
	}

}

func Test_advance(t *testing.T) {
	input := "123"
	p := setupParserForTest(input)

	p.advance()

	if string(p.previous.Lexeme) != "" {
		t.Error("Expected previous token Lexeme to be ''/(blank).")
	}

	if string(p.current.Lexeme) != input {
		t.Errorf("Expected current token Lexeme to be '%v'.", input)
	}

	p.advance()

	if string(p.previous.Lexeme) != input {
		t.Errorf("Expected previous token Lexeme to be '%v'.", input)
	}
}

func Test_consume(t *testing.T) {
	input := "123"
	msg := []byte("test consume error")
	p := setupParserForTest(input)

	p.current = p.lexer.ScanToken()

	p.consume(token.Number, msg)

	if string(p.previous.Lexeme) != input {
		t.Errorf("Expected current token Lexeme to be '%v'.", input)
	}

	p = setupParserForTest(input)

	p.current = p.lexer.ScanToken()

	p.consume(token.Nil, msg)

	if p.hadError != true {
		t.Error("Expected error to trigger for non matching token type.")
	}

}

func Test_beginScope(t *testing.T) {
	p := setupParserForTest("")

	p.beginScope()

	if p.compiler.scopeDepth != 1 {
		t.Errorf("Expected scopeDepth of 1, got %v.", p.compiler.scopeDepth)
	}
}

func Test_endScope(t *testing.T) {
	p := setupParserForTest("")

	p.compiler.scopeDepth = 1

	p.compiler.locals = append(p.compiler.locals,
		Local{depth: 1, name: token.Token{Type: token.Identifier, Lexeme: []byte("var1")}},
		Local{depth: 1, name: token.Token{Type: token.Identifier, Lexeme: []byte("var2")}},
		Local{depth: 1, name: token.Token{Type: token.Identifier, Lexeme: []byte("var3")}},
	)

	p.compiler.localCount = 3

	p.endScope()

	if p.compiler.scopeDepth != 0 {
		t.Errorf("Expected scopeDepth of 0, got %v.", p.compiler.scopeDepth)
	}

	if p.compiler.localCount != 1 {
		t.Errorf("Expected localCount of 1, got %v.", p.compiler.localCount)
	}
}

func Test_emitConstant(t *testing.T) {
	p := setupParserForTest("")

	con := value.NumberVal(123)

	p.emitConstant(con)

	checkConstants(t, p.chunk.Constants, []value.Value{con})
}

func Test_emitLoop(t *testing.T) {
	p := setupParserForTest("")
	p.chunk.Code = append(p.chunk.Code, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	loopStart := 5

	expected := p.chunk.Count() - loopStart + 2

	p.emitLoop(loopStart)

	result := int(p.chunk.Code[13]) << 8
	result |= int(p.chunk.Code[14] - 1)

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func Test_emitJump(t *testing.T) {
	p := setupParserForTest("")
	instruction := byte(opcode.JumpIfFalse)

	result := p.emitJump(instruction)

	opResult := p.chunk.Code[0]

	if opResult != instruction {
		t.Errorf("Expected %v, got %v", opcode.Name[instruction], opcode.Name[opResult])
	}

	expected := p.chunk.Count() - 2
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func Test_patchJump(t *testing.T) {
	p := setupParserForTest("")
	p.chunk.Code = append(p.chunk.Code, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	offset := 4

	p.patchJump(offset)

	expected := p.chunk.Count() - offset - 2
	result := int(p.chunk.Code[5])
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func Test_emitByte(t *testing.T) {
	p := setupParserForTest("")

	op := opcode.Add

	p.emitByte(op)

	checkOpcodes(t, p.chunk.Code, []byte{op})
}

func setupParserForTest(source string) *Parser {
	s := []byte(source)
	l := lexer.NewLexer(&s)
	co := NewCompiler()
	ch := chunk.NewChunk()
	return NewParser(l, ch, co)
}

func checkOpcodes(t *testing.T, actual []byte, expected []byte) {
	t.Helper()

	gotLen := len(actual)
	expectedLen := len(expected)

	if gotLen != expectedLen {
		t.Fatalf("Expected byte code slice with length '%v', got '%v'.", expectedLen, gotLen)
	}

	for i := 0; i < gotLen; i++ {
		actName, actExists := opcode.Name[actual[i]]
		expName, expExists := opcode.Name[expected[i]]
		if !actExists || !expExists {
			if !actExists {
				t.Errorf("Unknown opcode '%v' at code index %v of actual slice.", actual[i], i)
			}
			if !expExists {
				t.Errorf("Unknown opcode '%v' at code index %v of expected slice.", expected[i], i)
			}
		} else {
			if actual[i] != expected[i] {
				t.Errorf("Expected opcode '%v' at code index %v, got '%v'.", expName, i, actName)
			} else {
				switch expected[i] {
				case opcode.Constant, opcode.GetLocal, opcode.SetLocal,
					opcode.GetGlobal, opcode.DefineGlobal, opcode.SetGlobal:
					i++
					if actual[i] != expected[i] {
						t.Errorf("Expected %v with value %v at code index %v, got value %v.", expName, expected[i], i, actual[i])
					}
				case opcode.Loop, opcode.Jump, opcode.JumpIfFalse:
					actIndex := int(actual[i+1]) << 8
					actIndex |= int(actual[i+2])
					expIndex := int(expected[i+1]) << 8
					expIndex |= int(expected[i+2])
					i += 2
					if actIndex != expIndex {
						t.Errorf("Expected %v with value %v at code index %v, got value %v.", expName, expIndex, i, actIndex)
					}
				}
			}
		}
	}
}

func checkConstants(t *testing.T, actual []value.Value, expected []value.Value) {
	t.Helper()

	gotLen := len(actual)
	expectedLen := len(expected)

	if gotLen != expectedLen {
		t.Fatalf("Expected constants slice with length '%v', got '%v'.", expectedLen, gotLen)
	}

	for i := 0; i < gotLen; i++ {
		constant := actual[i]
		expectedConstant := expected[i]
		if constant != expectedConstant {
			t.Errorf("Expected constant '%v' at index %v, got '%v'.", expectedConstant, i, constant)
		}
	}
}
