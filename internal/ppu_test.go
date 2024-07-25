package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type tile struct {
	data [16]byte
}

func (t *tile) getPixels() [8][8]byte {
	colorValues := [8][8]byte{}

	for i := 0; i < 8; i++ {
		start := i * 2
		byte1 := t.data[start]
		byte2 := t.data[start+1]

		for b := byte(0); b < 8; b++ {
			mask := byte(0b0000_0001)

			bit1 := (byte1 >> (7 - b)) & mask
			bit2 := (byte2 >> (7 - b)) & mask

			colorValue := (bit2 << 1) | bit1
			colorValues[i][b] = colorValue
		}
	}

	return colorValues
}

func TestTileComputation(t *testing.T) {
	tl := tile{[16]byte{0x3C, 0x7E, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7E, 0x5E, 0x7E, 0x0A, 0x7C, 0x56, 0x38, 0x7C}}

	result := tl.getPixels()

	colorValues := [8][8]byte{
		{0, 2, 3, 3, 3, 3, 2, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 1, 3, 3, 3, 3, 0},
		{0, 1, 1, 1, 3, 1, 3, 0},
		{0, 3, 1, 3, 1, 3, 2, 0},
		{0, 2, 3, 3, 3, 2, 0, 0},
	}

	assert.Equal(t, colorValues, result)
}
