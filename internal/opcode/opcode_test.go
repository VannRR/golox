package opcode_test

import (
	"golox/internal/opcode"
	"testing"
)

func Test_Name(t *testing.T) {
	var i byte = 0
	for ; i < opcode.Return; i++ {
		if opcode.Name(i) == "OpUnknown" {
			t.Errorf("Opcode '%v' is unknown.", i)
		}
	}
}
