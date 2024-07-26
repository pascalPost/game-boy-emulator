package alu

import (
	"fmt"
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions/utils"
)

func BitwiseOrRegister(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr, register uint8, registerName string) {
	result := *registerA | register
	*registerA = result

	flags.ClearAll()
	if result == 0 {
		flags.SetZ()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("OR A, %s", registerName)
	description := "Performs a bitwise OR operation between the 8-bit A register and the 8-bit register r, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func BitwiseXorRegister(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr, register uint8, registerName string) {
	result := *registerA ^ register
	*registerA = result

	flags.ClearAll()
	if result == 0 {
		flags.SetZ()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("XOR A, %s", registerName)
	description := "XOR r: Performs a bitwise XOR operation between the 8-bit A register and the 8-bit register r, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func halfCarryAdd(a, b uint8) bool {
	return (a&0b1111)+(b&0b1111) > 0b1111
}

func carryAdd(a, b uint8) bool {
	return uint16(a)+uint16(b) > 0b1111_1111
}

func addImpl(registerA *uint8, n8 uint8, flags cpu.FlagsPtr) {
	a := *registerA
	result := a + n8
	*registerA = result

	flags.ClearAll()

	if result == 0 {
		flags.SetZ()
	}
	flags.SetN()
	if halfCarryAdd(a, n8) {
		flags.SetH()
	}
	if carryAdd(a, n8) {
		flags.SetC()
	}
}

func AddRegister(memory *cpu.Memory, programCounter uint16, registerA *uint8, register uint8, flags cpu.FlagsPtr, registerName string) {
	addImpl(registerA, register, flags)
	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("ADD %s", registerName)
	description := "Adds to the 8-bit A register, the 8-bit register r, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func AddIndirectHL(memory *cpu.Memory, programCounter uint16, registerA *uint8, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	addImpl(registerA, n8, flags)

	instructionLengthInBytes := 1
	instruction := "ADD [HL]"
	description := "Adds from the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func AddImmediate(memory *cpu.Memory, programCounter *uint16, registerA *uint8, flags cpu.FlagsPtr) {
	pc := *programCounter
	n8 := utils.ReadUnsigned8(memory, programCounter)
	addImpl(registerA, n8, flags)
	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("ADD 0x%02X", n8)
	description := "Adds from the 8-bit A register, the immediate data n, and stores the result back into the A register."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func halfCarrySub(a, b uint8) bool {
	return (a & 0b1111) < (b & 0b1111)
}

func carrySub(a, b uint8) bool {
	return a < b
}

func subtractImpl(registerA uint8, n8 uint8, flags cpu.FlagsPtr) uint8 {
	result := registerA - n8

	flags.ClearAll()

	if result == 0 {
		flags.SetZ()
	}
	flags.SetN()
	if halfCarrySub(registerA, n8) {
		flags.SetH()
	}
	if carrySub(registerA, n8) {
		flags.SetC()
	}

	return result
}

func SubtractRegister(memory *cpu.Memory, programCounter uint16, registerA *uint8, register uint8, flags cpu.FlagsPtr, registerName string) {
	res := subtractImpl(*registerA, register, flags)
	*registerA = res
	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("SUB %s", registerName)
	description := "Subtracts from the 8-bit A register, the 8-bit register r, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func SubtractIndirectHL(memory *cpu.Memory, programCounter uint16, registerA *uint8, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	res := subtractImpl(*registerA, n8, flags)
	*registerA = res

	instructionLengthInBytes := 1
	instruction := "SUB [HL]"
	description := "Subtracts from the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL, and stores the result back into the A register."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func SubtractImmediate(memory *cpu.Memory, programCounter *uint16, registerA *uint8, flags cpu.FlagsPtr) {
	pc := *programCounter

	n8 := utils.ReadUnsigned8(memory, programCounter)
	res := subtractImpl(*registerA, n8, flags)
	*registerA = res

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("SUB 0x%02X", n8)
	description := "Subtracts from the 8-bit A register, the immediate data n, and stores the result back into the A register."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func IncrementRegister(memory *cpu.Memory, programCounter uint16, register *uint8, flags cpu.FlagsPtr, registerName string) {
	old := *register
	*register = old + 1

	if *register == 0 {
		flags.SetZ()
	} else {
		flags.ClearZ()
	}
	flags.ClearN()
	if halfCarryAdd(old, 1) {
		flags.SetH()
	} else {
		flags.ClearH()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("INC %s", registerName)
	description := "INC rr: Increments data in the 16-bit register rr."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func Increment16BitRegister(memory *cpu.Memory, programCounter uint16, register *uint16, registerName string) {
	*register++

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("INC %s", registerName)
	description := "INC rr: Increments data in the 16-bit register rr."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func Decrement16BitRegister(memory *cpu.Memory, programCounter uint16, register *uint16, registerName string) {
	*register--

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("DEC %s", registerName)
	description := "DEC rr: Decrements data in the 16-bit register rr."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func CompareImmediate(memory *cpu.Memory, programCounter *uint16, registerA uint8, flags cpu.FlagsPtr) {
	pc := *programCounter

	n8 := utils.ReadUnsigned8(memory, programCounter)
	subtractImpl(registerA, n8, flags)

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("CP 0x%02X", n8)
	description := "Subtracts from the 8-bit A register, the immediate data n, and updates ﬂags based on the result. This instruction is basically identical to SUB n, but does not update the A register."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}
