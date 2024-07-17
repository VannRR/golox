package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const DEBUG_TRACE_EXECUTION bool = true

func main() {
	vm := NewVM()

	if argc := len(os.Args); argc == 1 {
		repl(vm)
	} else if argc == 2 {
		runFile(vm, os.Args[1])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: golox [path]\n")
		os.Exit(64)
	}

	vm.Free()
}

func repl(vm *VM) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")

		if line, _ := reader.ReadBytes('\n'); line[0] != '\n' {
			vm.Interpret(&line)
		}
	}
}

func runFile(vm *VM, path string) {
	source := readFile(path)
	result := vm.Interpret(source)
	*source = nil

	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}
}

func readFile(path string) *[]byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 0)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		buf = append(buf, line...)
	}

	return &buf
}
