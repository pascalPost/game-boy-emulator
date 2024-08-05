package ppu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func flatten(slice [][]byte) []byte {
	l := 0
	for _, v := range slice {
		l += len(v)
	}
	data := make([]byte, 0, l)
	for _, v := range slice {
		data = append(data, v[:]...)
	}
	return data
}

var tile = []byte{0x3C, 0x7E, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7E, 0x5E, 0x7E, 0x0A, 0x7C, 0x56, 0x38, 0x7C}

var colorValues = flatten([][]byte{
	{0, 2, 3, 3, 3, 3, 2, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 1, 3, 3, 3, 3, 0},
	{0, 1, 1, 1, 3, 1, 3, 0},
	{0, 3, 1, 3, 1, 3, 2, 0},
	{0, 2, 3, 3, 3, 2, 0, 0},
})

func TestTileComputation(t *testing.T) {
	result := ComputePixelColors(tile)
	assert.Equal(t, colorValues, result)

}

//func TestTilePlotting(t *testing.T) {
//	PlotTile(colorValues)
//}
