package internal

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestDisassembler(t *testing.T) {
	data := []struct {
		code        []byte
		instruction string
	}{
		{[]byte{0x31, 0x00, 0xE0}, "LD SP 0xE000"},
		{[]byte{0xCD, 0xA3, 0x17}, "CALL 0x17A3"},
		{[]byte{0xC3, 0xB6, 0x15}, "JP 0x15B6"},
		{[]byte{0xF5}, "PUSH AF"},
		{[]byte{0x3E, 0x01}, "LD A 0x01"},
		{[]byte{0xEA, 0x1C, 0xC3}, "LD 0xC31C A"},
		{[]byte{0xF1}, "POP AF"},
	}

	opcodes, err := ParseOpcodes()
	assert.NoError(t, err)

	for _, d := range data {
		instructions := Disassemble(d.code, 0, opcodes)
		assert.Equal(t, 1, len(instructions))
		assert.Equal(t, d.instruction, instructions[0].Line)
	}
}

func TestDisassembleSnake(t *testing.T) {
	snakeUrl := "https://hh3.gbdev.io/static/database-gb/entries/snake-gb/snake.gb"
	resp, err := http.Get(snakeUrl)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	}()

	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	opcodes, err := ParseOpcodes()
	assert.NoError(t, err)

	instructions := Disassemble(data, 0x0150, opcodes)
	assert.Greater(t, len(instructions), 0)
}
