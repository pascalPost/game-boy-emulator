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

func TestFlagStatusSet(t *testing.T) {
	var num uint8 = 0b1111_0000

	assert.True(t, isBit7Set(num))
	assert.True(t, isBit6Set(num))
	assert.True(t, isBit5Set(num))
	assert.True(t, isBit4Set(num))
}

func TestFlagStatusNotSet(t *testing.T) {
	var num uint8 = 0b0000_0000

	assert.False(t, isBit7Set(num))
	assert.False(t, isBit6Set(num))
	assert.False(t, isBit5Set(num))
	assert.False(t, isBit4Set(num))
}

func TestSetAndClearBits(t *testing.T) {
	tests := []struct {
		init   uint8
		f      func(value *uint8)
		result uint8
	}{
		{0b0000_0000, setBit7, 0b1000_0000},
		{0b1111_1111, setBit7, 0b1111_1111},
		{0b0000_0000, setBit6, 0b0100_0000},
		{0b1111_1111, setBit6, 0b1111_1111},
		{0b0000_0000, setBit5, 0b0010_0000},
		{0b1111_1111, setBit5, 0b1111_1111},
		{0b0000_0000, setBit4, 0b0001_0000},
		{0b1111_1111, setBit4, 0b1111_1111},

		{0b0000_0000, clearBit7, 0b0000_0000},
		{0b1111_1111, clearBit7, 0b0111_1111},
		{0b0000_0000, clearBit6, 0b0000_0000},
		{0b1111_1111, clearBit6, 0b1011_1111},
		{0b0000_0000, clearBit5, 0b0000_0000},
		{0b1111_1111, clearBit5, 0b1101_1111},
		{0b0000_0000, clearBit4, 0b0000_0000},
		{0b1111_1111, clearBit4, 0b1110_1111},

		{0b0000_0000, clearBits7to4, 0b0000_0000},
		{0b1111_1111, clearBits7to4, 0b0000_1111},
	}

	for _, test := range tests {
		value := test.init
		test.f(&value)
		assert.Equal(t, test.result, value)
	}
}
