package utils

import "github.com/pascalPost/game-boy-emulator/internal/cpu"

func ReadUnsigned8(memory *cpu.Memory, programCounter *uint16) uint8 {
	v := memory.Read(*programCounter)
	*programCounter++
	return v
}

func ReadUnsigned16(memory *cpu.Memory, programCounter *uint16) uint16 {
	nnLSB := ReadUnsigned8(memory, programCounter)
	nnMSB := ReadUnsigned8(memory, programCounter)
	nn := Unsigned16(nnLSB, nnMSB)
	return nn
}

func Unsigned16(leastSignificantByte uint8, mostSignificantByte uint8) uint16 {
	// littleEndian
	return uint16(leastSignificantByte) | uint16(mostSignificantByte)<<8
}
