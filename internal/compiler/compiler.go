package compiler

import (
	"fmt"
	"golox/internal/lexer"
	"golox/internal/token"
)

func Compile(source *[]byte) {
	l := lexer.NewLexer(source)
	line := -1
	for {
		t := l.ScanToken()
		if t.Line != line {
			fmt.Printf("%4d ", t.Line)
			line = t.Line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", t.Type, t.Lexeme)

		if t.Type == token.Eof {
			break
		}
	}
}
