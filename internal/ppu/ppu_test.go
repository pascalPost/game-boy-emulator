package ppu

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestTileComputation(t *testing.T) {
	tile := [16]byte{0x3C, 0x7E, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7E, 0x5E, 0x7E, 0x0A, 0x7C, 0x56, 0x38, 0x7C}

	result := GetPixels(tile[:])

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

	slices.Reverse(colorValues[:])

	l := 0
	for _, v := range colorValues {
		l += len(v)
	}
	colors := make([]byte, 0, l)
	for _, v := range colorValues {
		colors = append(colors, v[:]...)
	}

	//plotTile(colors)
}
