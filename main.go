package main

const DEBUG_TRACE_EXECUTION bool = true

func main() {
    vm := NewVM()

    var chunk = NewChunk()
	chunk.WriteConstant(1.2, 123)
	chunk.Write(OP_RETURN, 123)
	//chunk.Disassemble("test chunk")
    vm.Interpret(chunk)

    vm.Free()
    chunk.Free()
}
