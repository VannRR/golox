package main

import "fmt"

func Compile(source *[]byte) {
	scanner := NewScanner(source)
	line := -1
	for {
		token := scanner.ScanToken()
		if token.line != line {
			fmt.Printf("%4d ", token.line)
			line = token.line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("%2d '%s'\n", token.tokenType, token.lexeme)

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}
