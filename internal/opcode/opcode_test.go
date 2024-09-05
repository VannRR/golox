package opcode_test

import (
	"golox/internal/opcode"
	"testing"
)

func Test_Name(t *testing.T) {
	for i := byte(0); i < opcode.Return; i++ {
		_, exists := opcode.Name[i]
		if !exists {
			t.Errorf("Opcode '%v' is unknown.", i)
		}
	}
}
