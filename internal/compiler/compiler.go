package compiler

import (
	"bytes"
	"fmt"
	"github.com/VannRR/golox/internal/chunk"
	"github.com/VannRR/golox/internal/common"
	"github.com/VannRR/golox/internal/debug"
	"github.com/VannRR/golox/internal/lexer"
	"github.com/VannRR/golox/internal/object"
	"github.com/VannRR/golox/internal/opcode"
	"github.com/VannRR/golox/internal/token"
	"github.com/VannRR/golox/internal/value"
	"os"
	"strconv"
)

const (
	PrecNone       Precedence = iota
	PrecAssignment            // =
	PrecOr                    // or
	PrecAnd                   // and
	PrecEquality              // == !=
	PrecComparison            // < > <= >=
	PrecTerm                  // + -
	PrecFactor                // * /
	PrecUnary                 // ! -
	PrecCall                  // . ()
	PrecPrimary
)

var rules [token.Eof + 1]ParseRule

func init() {
	rules[token.LeftParen] = ParseRule{(*Parser).grouping, nil, PrecNone}
	rules[token.RightParen] = ParseRule{nil, nil, PrecNone}
	rules[token.LeftBrace] = ParseRule{nil, nil, PrecNone}
	rules[token.RightBrace] = ParseRule{nil, nil, PrecNone}
	rules[token.Comma] = ParseRule{nil, nil, PrecNone}
	rules[token.Dot] = ParseRule{nil, nil, PrecNone}
	rules[token.Minus] = ParseRule{(*Parser).unary, (*Parser).binary, PrecTerm}
	rules[token.Plus] = ParseRule{nil, (*Parser).binary, PrecTerm}
	rules[token.Semicolon] = ParseRule{nil, nil, PrecNone}
	rules[token.Slash] = ParseRule{nil, (*Parser).binary, PrecFactor}
	rules[token.Star] = ParseRule{nil, (*Parser).binary, PrecFactor}
	rules[token.Percent] = ParseRule{nil, (*Parser).binary, PrecFactor}
	rules[token.Bang] = ParseRule{(*Parser).unary, nil, PrecNone}
	rules[token.BangEqual] = ParseRule{nil, (*Parser).binary, PrecEquality}
	rules[token.Equal] = ParseRule{nil, nil, PrecNone}
	rules[token.EqualEqual] = ParseRule{nil, (*Parser).binary, PrecEquality}
	rules[token.Greater] = ParseRule{nil, (*Parser).binary, PrecComparison}
	rules[token.GreaterEqual] = ParseRule{nil, (*Parser).binary, PrecComparison}
	rules[token.Less] = ParseRule{nil, (*Parser).binary, PrecComparison}
	rules[token.LessEqual] = ParseRule{nil, (*Parser).binary, PrecComparison}
	rules[token.Identifier] = ParseRule{(*Parser).variable, nil, PrecNone}
	rules[token.String] = ParseRule{(*Parser).string, nil, PrecNone}
	rules[token.Number] = ParseRule{(*Parser).number, nil, PrecNone}
	rules[token.And] = ParseRule{nil, (*Parser).and, PrecAnd}
	rules[token.Class] = ParseRule{nil, nil, PrecNone}
	rules[token.Else] = ParseRule{nil, nil, PrecNone}
	rules[token.False] = ParseRule{(*Parser).literal, nil, PrecNone}
	rules[token.For] = ParseRule{nil, nil, PrecNone}
	rules[token.Fun] = ParseRule{nil, nil, PrecNone}
	rules[token.If] = ParseRule{nil, nil, PrecNone}
	rules[token.Nil] = ParseRule{(*Parser).literal, nil, PrecNone}
	rules[token.Or] = ParseRule{nil, (*Parser).or, PrecOr}
	rules[token.Print] = ParseRule{nil, nil, PrecNone}
	rules[token.Return] = ParseRule{nil, nil, PrecNone}
	rules[token.Super] = ParseRule{nil, nil, PrecNone}
	rules[token.This] = ParseRule{nil, nil, PrecNone}
	rules[token.True] = ParseRule{(*Parser).literal, nil, PrecNone}
	rules[token.Var] = ParseRule{nil, nil, PrecNone}
	rules[token.While] = ParseRule{nil, nil, PrecNone}
	rules[token.Error] = ParseRule{nil, nil, PrecNone}
	rules[token.Eof] = ParseRule{nil, nil, PrecNone}
}

type ParseFn = func(*Parser, bool)
type Precedence = uint8

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

func getRule(tt token.TokenType) *ParseRule {
	return &rules[tt]
}

type Local struct {
	name  token.Token
	depth int
}

type Compiler struct {
	locals     []Local
	localCount int
	scopeDepth int
}

func NewCompiler() *Compiler {
	return &Compiler{
		locals:     make([]Local, 1),
		localCount: 0,
		scopeDepth: 0,
	}
}

type Parser struct {
	lexer     *lexer.Lexer
	chunk     *chunk.Chunk
	compiler  *Compiler
	current   token.Token
	previous  token.Token
	hadError  bool
	panicMode bool
}

func NewParser(l *lexer.Lexer, ch *chunk.Chunk, co *Compiler) *Parser {
	return &Parser{
		lexer:     l,
		chunk:     ch,
		compiler:  co,
		current:   token.Token{},
		previous:  token.Token{},
		hadError:  false,
		panicMode: false,
	}
}

func Compile(source *[]byte, ch *chunk.Chunk) bool {
	l := lexer.NewLexer(source)
	co := NewCompiler()
	p := NewParser(l, ch, co)
	p.advance()
	for !p.match(token.Eof) {
		p.declaration()
	}
	p.endCompiler()
	return !p.hadError
}

func (p *Parser) declaration() {
	if p.match(token.Var) {
		p.varDeclaration()
	} else {
		p.statement()
	}

	if p.panicMode {
		p.synchronize()
	}
}

func (p *Parser) statement() {
	if p.match(token.Print) {
		p.printStatement()
	} else if p.match(token.For) {
		p.forStatement()
	} else if p.match(token.If) {
		p.ifStatement()
	} else if p.match(token.While) {
		p.whileStatement()
	} else if p.match(token.LeftBrace) {
		p.beginScope()
		p.block()
		p.endScope()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) printStatement() {
	p.expression()
	p.consume(token.Semicolon, []byte("Expect ';' after value."))
	p.emitByte(opcode.Print)
}

func (p *Parser) forStatement() {
	p.beginScope()
	p.consume(token.LeftParen, []byte("Expect '(' after 'for'."))
	if p.match(token.Semicolon) {
		// No initializer.
	} else if p.match(token.Var) {
		p.varDeclaration()
	} else {
		p.expressionStatement()
	}

	loopStart := p.chunk.Count()
	exitJump := -1
	if !p.match(token.Semicolon) {
		p.expression()
		p.consume(token.Semicolon, []byte("Expect ';' after loop condition."))

		exitJump = p.emitJump(opcode.JumpIfFalse)
		p.emitByte(opcode.Pop)
	}

	if !p.match(token.RightParen) {
		bodyJump := p.emitJump(opcode.Jump)
		incrementStart := p.chunk.Count()
		p.expression()
		p.emitByte(opcode.Pop)
		p.consume(token.RightParen, []byte("Expect ')' after for clauses."))

		p.emitLoop(loopStart)
		loopStart = incrementStart
		p.patchJump(bodyJump)
	}

	p.statement()
	p.emitLoop(loopStart)

	if exitJump != -1 {
		p.patchJump(exitJump)
		p.emitByte(opcode.Pop)
	}

	p.endScope()
}

func (p *Parser) ifStatement() {
	p.consume(token.LeftParen, []byte("Expect '(' after 'if'."))
	p.expression()
	p.consume(token.RightParen, []byte("Expect ')' after condition."))

	thenJump := p.emitJump(opcode.JumpIfFalse)
	p.emitByte(opcode.Pop)
	p.statement()

	elseJump := p.emitJump(opcode.Jump)

	p.patchJump(thenJump)
	p.emitByte(opcode.Pop)

	if p.match(token.Else) {
		p.statement()
	}
	p.patchJump(elseJump)
}

func (p *Parser) whileStatement() {
	loopStart := p.chunk.Count()
	p.consume(token.LeftParen, []byte("Expect '(' after 'while'."))
	p.expression()
	p.consume(token.RightParen, []byte("Expect ')' after condition."))

	exitJump := p.emitJump(opcode.JumpIfFalse)
	p.emitByte(opcode.Pop)
	p.statement()
	p.emitLoop(loopStart)

	p.patchJump(exitJump)
	p.emitByte(opcode.Pop)
}

func (p *Parser) expressionStatement() {
	p.expression()
	p.consume(token.Semicolon, []byte("Expect ';' after expression."))
	p.emitByte(opcode.Pop)
}

func (p *Parser) synchronize() {
	p.panicMode = false

	for p.current.Type != token.Eof {
		if p.previous.Type == token.Semicolon {
			return
		}
		switch p.current.Type {
		case token.Class, token.Fun, token.Var, token.For,
			token.If, token.While, token.Print, token.Return:
			return
		default:
			// Do nothing
		}
		p.advance()
	}
}

func (p *Parser) varDeclaration() {
	global := p.parseVariable([]byte("Expect variable name."))

	if p.match(token.Equal) {
		p.expression()
	} else {
		p.emitByte(opcode.Nil)
	}
	p.consume(token.Semicolon,
		[]byte("Expect ';' after variable declaration."))

	p.defineVariable(global)
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecAssignment)
}

func (p *Parser) block() {
	for !p.check(token.RightBrace) && !p.check(token.Eof) {
		p.declaration()
	}
	p.consume(token.RightBrace, []byte("Expect '}' after block."))
}

func (p *Parser) grouping(canAssign bool) {
	p.expression()
	p.consume(token.RightParen, []byte("Expect ')' after expression."))
}

func (p *Parser) number(canAssign bool) {
	v, err := strconv.ParseFloat(string(p.previous.Lexeme), 64)
	if err != nil {
		panic(err)
	}
	p.emitConstant(value.NumberVal(v))
}

func (p *Parser) or(canAssign bool) {
	elseJump := p.emitJump(opcode.JumpIfFalse)
	endJump := p.emitJump(opcode.Jump)

	p.patchJump(elseJump)
	p.emitByte(opcode.Pop)

	p.parsePrecedence(PrecOr)
	p.patchJump(endJump)
}

func (p *Parser) variable(canAssign bool) {
	p.namedVariable(p.previous, canAssign)
}

func (p *Parser) namedVariable(name token.Token, canAssign bool) {
	var getOp, setOp uint8
	index := p.resolveLocal(&name)
	if index != -1 {
		getOp = opcode.GetLocal
		setOp = opcode.SetLocal
	} else {
		index = p.identifierConstant(&name)
		getOp = opcode.GetGlobal
		setOp = opcode.SetGlobal
	}

	if canAssign && p.match(token.Equal) {
		p.expression()
		p.chunk.WriteIndexWithCheck(index, setOp, p.previous.Line)
	} else {
		p.chunk.WriteIndexWithCheck(index, getOp, p.previous.Line)
	}
}

func (p *Parser) resolveLocal(name *token.Token) int {
	for i := p.compiler.localCount - 1; i >= 0; i-- {
		local := &p.compiler.locals[i]
		if identifiersEqual(name, &local.name) {
			if local.depth == -1 {
				p.error([]byte("Can't read local variable in its own initializer."))
			}
			return int(i)
		}
	}
	return -1
}

func (p *Parser) string(canAssign bool) {
	p.emitConstant(object.ObjString(string(p.previous.Lexeme)[1 : len(p.previous.Lexeme)-1]))
}

func (p *Parser) binary(canAssign bool) {
	operatorType := p.previous.Type
	rule := getRule(operatorType)
	p.parsePrecedence(Precedence(rule.precedence + 1))

	switch operatorType {
	case token.BangEqual:
		p.emitByte(opcode.NotEqual)
	case token.EqualEqual:
		p.emitByte(opcode.Equal)
	case token.Greater:
		p.emitByte(opcode.Greater)
	case token.GreaterEqual:
		p.emitByte(opcode.GreaterEqual)
	case token.Less:
		p.emitByte(opcode.Less)
	case token.LessEqual:
		p.emitByte(opcode.LessEqual)
	case token.Plus:
		p.emitByte(opcode.Add)
	case token.Minus:
		p.emitByte(opcode.Subtract)
	case token.Star:
		p.emitByte(opcode.Multiply)
	case token.Slash:
		p.emitByte(opcode.Divide)
	case token.Percent:
		p.emitByte(opcode.Modulo)
	default:
		panic("binary parser, unknown operator type")
	}
}

func (p *Parser) literal(canAssign bool) {
	switch p.previous.Type {
	case token.False:
		p.emitByte(opcode.False)
	case token.Nil:
		p.emitByte(opcode.Nil)
	case token.True:
		p.emitByte(opcode.True)
	default:
		panic("literal parser, unknown operator type")
	}
}

func (p *Parser) unary(canAssign bool) {
	operatorType := p.previous.Type

	p.parsePrecedence(PrecUnary)

	switch operatorType {
	case token.Bang:
		p.emitByte(opcode.Not)
	case token.Minus:
		p.emitByte(opcode.Negate)
	default:
		panic("unary parser, unknown operator type")
	}
}

func (p *Parser) parsePrecedence(precedence Precedence) {
	p.advance()
	prefixRule := getRule(p.previous.Type).prefix

	if prefixRule == nil {
		p.error([]byte("Expect expression."))
		return
	}

	canAssign := precedence <= PrecAssignment
	prefixRule(p, canAssign)

	for precedence <= getRule(p.current.Type).precedence {
		p.advance()
		infixRule := getRule(p.previous.Type).infix
		infixRule(p, canAssign)
	}

	if canAssign && p.match(token.Equal) {
		p.error([]byte("Invalid assignment target."))
	}
}

func (p *Parser) parseVariable(errorMessage []byte) int {
	p.consume(token.Identifier, errorMessage)

	p.declareVariable()
	if p.compiler.scopeDepth > 0 {
		return 0
	}

	return p.identifierConstant(&p.previous)
}

func (p *Parser) identifierConstant(name *token.Token) int {
	index := p.chunk.AddConstant(object.ObjString(string(name.Lexeme)))
	p.chunk.WriteIndexWithCheck(index, opcode.Constant, p.current.Line)
	return index
}

func (p *Parser) declareVariable() {
	if p.compiler.scopeDepth == 0 {
		return
	}

	name := &p.previous

	for i := p.compiler.localCount - 1; i >= 0; i-- {
		local := &p.compiler.locals[i]
		if local.depth != -1 && local.depth < p.compiler.scopeDepth {
			break
		}

		if identifiersEqual(name, &local.name) {
			p.error([]byte("Already a variable with this name in this scope."))
		}
	}

	p.addLocal(*name)
}

func identifiersEqual(a *token.Token, b *token.Token) bool {
	return bytes.Equal(a.Lexeme, b.Lexeme)
}

func (p *Parser) addLocal(name token.Token) {
	if p.compiler.localCount >= common.Uint24Max {
		p.error([]byte("Too many local variables in function."))
		return
	}

	local := &p.compiler.locals[p.compiler.localCount]
	p.compiler.localCount++
	local.name = name
	local.depth = -1
}

func (p *Parser) defineVariable(global int) {
	if p.compiler.scopeDepth > 0 {
		p.markInitialized()
		return
	}
	p.chunk.WriteIndexWithCheck(global, opcode.DefineGlobal, p.previous.Line)
}

func (p *Parser) and(canAssign bool) {
	endJump := p.emitJump(opcode.JumpIfFalse)

	p.emitByte(opcode.Pop)
	p.parsePrecedence(PrecAnd)

	p.patchJump(endJump)
}

func (p *Parser) markInitialized() {
	p.compiler.locals[p.compiler.localCount-1].depth = p.compiler.scopeDepth
}

func (p *Parser) error(message []byte) {
	p.errorAt(&p.previous, message)
}

func (p *Parser) errorAtCurrent(message []byte) {
	p.errorAt(&p.current, message)
}

func (p *Parser) errorAt(t *token.Token, message []byte) {
	if p.panicMode {
		return
	}
	p.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", t.Line)

	if t.Type == token.Eof {
		fmt.Fprintf(os.Stderr, " at end")
	} else if t.Type == token.Error {
		// Nothing.
	} else {
		fmt.Fprintf(os.Stderr, " at %s", string(t.Lexeme))
	}

	fmt.Fprintf(os.Stderr, ": %s\n", string(message))
	p.hadError = true
}

func (p *Parser) match(tt token.TokenType) bool {
	if !p.check(tt) {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) check(tt token.TokenType) bool {
	return p.current.Type == tt
}

func (p *Parser) advance() {
	p.previous = p.current

	for {
		p.current = p.lexer.ScanToken()
		if p.current.Type != token.Error {
			break
		}
		p.errorAtCurrent(p.current.Lexeme)
	}
}

func (p *Parser) consume(tt token.TokenType, message []byte) {
	if p.current.Type == tt {
		p.advance()
		return
	}

	p.errorAtCurrent(message)
}

func (p *Parser) endCompiler() {
	p.emitReturn()
	if debug.PrintCode && !p.hadError {
		debug.DisassembleChunk(p.chunk, "code")
	}
}

func (p *Parser) beginScope() {
	p.compiler.scopeDepth++
}

func (p *Parser) endScope() {
	p.compiler.scopeDepth--

	for p.compiler.localCount > 0 &&
		p.compiler.locals[p.compiler.localCount-1].depth >
			p.compiler.scopeDepth {
		p.emitByte(opcode.Pop)
		p.compiler.localCount--
	}
}

func (p *Parser) emitReturn() {
	p.emitByte(opcode.Return)
}

func (p *Parser) emitConstant(v value.Value) {
	index := p.chunk.AddConstant(v)
	p.chunk.WriteIndexWithCheck(index, opcode.Constant, p.previous.Line)
}

func (p *Parser) emitLoop(loopStart int) {
	p.emitByte(opcode.Loop)

	offset := p.chunk.Count() - loopStart + 2
	if offset > common.Uint16Max {
		p.error([]byte("Loop body too large."))
	}

	p.emitByte(byte(offset >> 8))
	p.emitByte(byte(offset))
}

func (p *Parser) emitJump(instruction byte) int {
	p.emitByte(instruction)
	p.emitByte(0xff)
	p.emitByte(0xff)
	return p.chunk.Count() - 2
}

func (p *Parser) patchJump(offset int) {
	jump := p.chunk.Count() - offset - 2

	if jump > common.Uint16Max {
		p.error([]byte("Too much code to jump over."))
	}

	p.chunk.Code[offset] = byte(jump >> 8)
	p.chunk.Code[offset+1] = byte(jump)
}

// func (p *Parser) emitBytes(byte1 byte, byte2 byte) {
// 	p.emitByte(byte1)
// 	p.emitByte(byte2)
// }

func (p *Parser) emitByte(byte byte) {
	p.chunk.Write(byte, p.previous.Line)
}
