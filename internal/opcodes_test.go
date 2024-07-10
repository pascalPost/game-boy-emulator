package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseOpcodes(t *testing.T) {
	codes, err := ParseOpcodes()
	assert.NoError(t, err)

	unprefixed_0x00 := Opcode{
		Mnemonic:  "NOP",
		Bytes:     1,
		Cycles:    []int{4},
		Operands:  []Operands{},
		Immediate: true,
		Flags: Flags{
			Z: "-",
			N: "-",
			H: "-",
			C: "-",
		},
	}

	assert.Equal(t, unprefixed_0x00, codes.UnPrefixed[ByteKey{0x00}])
}
