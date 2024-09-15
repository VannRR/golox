package debug_test

import (
	"bytes"
	"fmt"
	"github.com/VannRR/golox/internal/chunk"
	"github.com/VannRR/golox/internal/debug"
	"github.com/VannRR/golox/internal/object"
	"github.com/VannRR/golox/internal/opcode"
	"github.com/VannRR/golox/internal/value"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDisassembleInstruction(t *testing.T) {
	c := &chunk.Chunk{
		Code: []byte{
			opcode.Constant, 0,
			opcode.Add,
			opcode.Modulo,
		},
		Constants: []value.Value{object.ObjString("Hello, World!")},
	}

	tests := []struct {
		offset int
		want   []string
	}{
		{0, []string{"0000", "0", "OpConstant", "0", "'Hello, World!'"}},
		{2, []string{"0002", "|", "OpAdd"}},
		{3, []string{"0003", "|", "OpModulo"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Offset %d", tt.offset), func(t *testing.T) {
			output := captureOutput(func() {
				debug.DisassembleInstruction(c, tt.offset)
			})

			for _, part := range tt.want {
				if !strings.Contains(output, part) {
					t.Errorf("Expected output to contain: %s\nGot:\n%s", part, output)
				}
			}
		})
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
