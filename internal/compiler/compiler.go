package compiler

import (
	"fmt"
	"golox/internal/chunk"
	"golox/internal/debug"
	"golox/internal/lexer"
	"golox/internal/opcode"
	"golox/internal/token"
	"golox/internal/value"
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
	rules[token.Identifier] = ParseRule{nil, nil, PrecNone}
	rules[token.String] = ParseRule{(*Parser).string, nil, PrecNone}
	rules[token.Number] = ParseRule{(*Parser).number, nil, PrecNone}
	rules[token.And] = ParseRule{nil, nil, PrecNone}
	rules[token.Class] = ParseRule{nil, nil, PrecNone}
	rules[token.Else] = ParseRule{nil, nil, PrecNone}
	rules[token.False] = ParseRule{(*Parser).literal, nil, PrecNone}
	rules[token.For] = ParseRule{nil, nil, PrecNone}
	rules[token.Fun] = ParseRule{nil, nil, PrecNone}
	rules[token.If] = ParseRule{nil, nil, PrecNone}
	rules[token.Nil] = ParseRule{(*Parser).literal, nil, PrecNone}
	rules[token.Or] = ParseRule{nil, nil, PrecNone}
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

type ParseFn = func(*Parser)
type Precedence = uint8

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

func getRule(tt token.Type) *ParseRule {
	return &rules[tt]
}

type Parser struct {
	lexer     lexer.Lexer
	chunk     *chunk.Chunk
	current   token.Token
	previous  token.Token
	hadError  bool
	panicMode bool
}

func NewParser(source *[]byte, c *chunk.Chunk) *Parser {
	return &Parser{
		lexer:     *lexer.NewLexer(source),
		chunk:     c,
		current:   token.Token{},
		previous:  token.Token{},
		hadError:  false,
		panicMode: false,
	}
}

func Compile(source *[]byte, c *chunk.Chunk) bool {
	p := NewParser(source, c)
	p.advance()
	p.expression()
	p.consume(token.Eof, []byte("Expect end of expression."))
	p.endCompiler()
	return !p.hadError
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecAssignment)
}

func (p *Parser) grouping() {
	p.expression()
	p.consume(token.RightParen, []byte("Expect ')' after expression."))
}

func (p *Parser) number() {
	v, err := strconv.ParseFloat(string(p.previous.Lexeme), 64)
	if err != nil {
		panic(err)
	}
	p.emitConstant(value.NewNumber(v))
}

func (p *Parser) string() {
	p.emitConstant(value.NewObjString(string(p.previous.Lexeme)[1 : len(p.previous.Lexeme)-1]))
}

func (p *Parser) binary() {
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

func (p *Parser) literal() {
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

func (p *Parser) unary() {
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

	prefixRule(p)

	for precedence <= getRule(p.current.Type).precedence {
		p.advance()
		infixRule := getRule(p.previous.Type).infix
		infixRule(p)
	}
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

func (p *Parser) consume(tt token.Type, message []byte) {
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
func (p *Parser) emitReturn() {
	p.emitByte(opcode.Return)
}

func (p *Parser) emitConstant(v value.Value) {
	errMsg, hasErr := p.chunk.WriteConstant(v, p.previous.Line)
	if hasErr {
		p.error([]byte(errMsg))
	}
}

func (p *Parser) emitBytes(byte1 byte, byte2 byte) {
	p.emitByte(byte1)
	p.emitByte(byte2)
}

func (p *Parser) emitByte(byte byte) {
	p.chunk.Write(byte, p.previous.Line)
}
