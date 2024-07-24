package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterSplit(t *testing.T) {
	var num uint16 = 0x1234

	assert.Equal(t, uint8(0x12), highPart(num))
	assert.Equal(t, uint8(0x34), lowPart(num))

	assert.Equal(t, uint8(0x12), *highPartPtr(&num))
	assert.Equal(t, uint8(0x34), *lowPartPtr(&num))
}

func TestIsBitSet(t *testing.T) {
	tests := []struct {
		num    uint8
		bit    []uint8
		result bool
	}{
		{0b1111_1111, []uint8{0, 1, 2, 3, 4, 5, 6, 7}, true},
		{0b0000_0000, []uint8{0, 1, 2, 3, 4, 5, 6, 7}, false},
	}

	for _, test := range tests {
		for _, bit := range test.bit {
			res, err := isBitSet(test.num, bit)
			assert.NoError(t, err)
			assert.Equal(t, test.result, res)
		}
	}

	_, err := isBitSet(0b1111_1111, 8)
	assert.Error(t, err)
}

func TestSetBit(t *testing.T) {
	tests := []struct {
		init   uint8
		bit    uint8
		result uint8
	}{
		{0b0000_0000, 7, 0b1000_0000},
		{0b1111_1111, 7, 0b1111_1111},
		{0b0000_0000, 6, 0b0100_0000},
		{0b1111_1111, 6, 0b1111_1111},
		{0b0000_0000, 5, 0b0010_0000},
		{0b1111_1111, 5, 0b1111_1111},
		{0b0000_0000, 4, 0b0001_0000},
		{0b1111_1111, 4, 0b1111_1111},
		{0b0000_0000, 3, 0b0000_1000},
		{0b1111_1111, 3, 0b1111_1111},
		{0b0000_0000, 2, 0b000_0100},
		{0b1111_1111, 2, 0b1111_1111},
		{0b0000_0000, 1, 0b0000_0010},
		{0b1111_1111, 1, 0b1111_1111},
		{0b0000_0000, 0, 0b0000_0001},
		{0b1111_1111, 0, 0b1111_1111},
	}

	for _, test := range tests {
		value := test.init
		err := setBit(&value, test.bit)
		assert.NoError(t, err)
		assert.Equal(t, test.result, value)
	}

	value := uint8(0)
	err := setBit(&value, 8)
	assert.Error(t, err)
}

func TestClearBit(t *testing.T) {
	tests := []struct {
		init   uint8
		bit    uint8
		result uint8
	}{
		{0b0000_0000, 7, 0b0000_0000},
		{0b1111_1111, 7, 0b0111_1111},
		{0b0000_0000, 6, 0b0000_0000},
		{0b1111_1111, 6, 0b1011_1111},
		{0b0000_0000, 5, 0b0000_0000},
		{0b1111_1111, 5, 0b1101_1111},
		{0b0000_0000, 4, 0b0000_0000},
		{0b1111_1111, 4, 0b1110_1111},
		{0b0000_0000, 3, 0b0000_0000},
		{0b1111_1111, 3, 0b1111_0111},
		{0b0000_0000, 2, 0b0000_0000},
		{0b1111_1111, 2, 0b1111_1011},
		{0b0000_0000, 1, 0b0000_0000},
		{0b1111_1111, 1, 0b1111_1101},
		{0b0000_0000, 0, 0b0000_0000},
		{0b1111_1111, 0, 0b1111_1110},
	}

	for _, test := range tests {
		value := test.init
		err := clearBit(&value, test.bit)
		assert.NoError(t, err)
		assert.Equal(t, test.result, value)
	}

	value := uint8(0)
	err := clearBit(&value, 8)
	assert.Error(t, err)

}

func TestClearBits7to4(t *testing.T) {
	tests := []struct {
		init   uint8
		result uint8
	}{
		{0b0000_0000, 0b0000_0000},
		{0b1111_1111, 0b0000_1111},
	}

	for _, test := range tests {
		value := test.init
		clearBits7to4(&value)
		assert.Equal(t, test.result, value)
	}
}
