package cpu

import (
	"errors"
	"unsafe"
)

func bitNumberError() error {
	return errors.New("A 8bit value has only 8 bits. I.e., the maximal bit to specify is bit number 7")
}

func IsBitSet(value uint8, bitNumber uint8) (bool, error) {
	if bitNumber > 7 {
		return false, bitNumberError()
	}
	mask := uint8(0b0000_0001) << bitNumber
	return (value & mask) != 0, nil
}

func setBit(value *uint8, bitNumber uint8) error {
	if bitNumber > 7 {
		return bitNumberError()
	}
	mask := uint8(0b0000_0001) << bitNumber
	*value = *value | mask
	return nil
}

func clearBit(value *uint8, bitNumber uint8) error {
	if bitNumber > 7 {
		return bitNumberError()
	}
	mask := uint8(0b1111_1111) ^ (uint8(0b0000_0001) << bitNumber)
	*value = *value & mask
	return nil
}

func clearBits7to4(value *uint8) {
	*value = *value & 0b0000_1111
}

func HighPart(value uint16) uint8 {
	return uint8(value >> 8)
}

func LowPart(value uint16) uint8 {
	return uint8(value & 0xff)
}

func HighPartPtr(ptr *uint16) *uint8 {
	// this might only work for little endian systems, if so swap with lowPartPtr
	uPtr := unsafe.Pointer(ptr)
	return (*uint8)(unsafe.Pointer(uintptr(uPtr) + 1))
}

func LowPartPtr(ptr *uint16) *uint8 {
	// this might only work for little endian systems, if so swap with highPartPtr
	uPtr := unsafe.Pointer(ptr)
	return (*uint8)(uPtr)
}

func checkBit(value uint8, bitNumber uint8) bool {
	res, err := IsBitSet(value, bitNumber)
	if err != nil {
		panic(err)
	}
	return res
}
