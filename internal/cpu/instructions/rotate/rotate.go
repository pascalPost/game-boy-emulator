package rotate

import (
	"fmt"
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
)

func leftImpl(n8 *uint8, flags cpu.FlagsPtr, setFlagZ bool) {
	isBit7Set, _ := cpu.IsBitSet(*n8, 7)

	*n8 <<= 1
	if flags.C() {
		*n8 += 1
	}

	flags.ClearAll()
	if setFlagZ && *n8 == 0 {
		flags.SetZ()
	}
	if isBit7Set {
		flags.SetC()
	}
}

func leftCircularImpl(n8 *uint8, flags cpu.FlagsPtr, setFlagZ bool) {
	isBit7Set, _ := cpu.IsBitSet(*n8, 7)

	*n8 <<= 1
	if isBit7Set {
		*n8 += 1
	}

	flags.ClearAll()
	if setFlagZ && *n8 == 0 {
		flags.SetZ()
	}
	if isBit7Set {
		flags.SetC()
	}
}

func rightImpl(n8 *uint8, flags cpu.FlagsPtr, setFlagZ bool) {
	isBit0Set, _ := cpu.IsBitSet(*n8, 0)

	*n8 >>= 1
	if flags.C() {
		*n8 |= 0b1000_0000
	}

	flags.ClearAll()
	if setFlagZ && *n8 == 0 {
		flags.SetZ()
	}
	if isBit0Set {
		flags.SetC()
	}
}

func rightCircularImpl(n8 *uint8, flags cpu.FlagsPtr, setFlagZ bool) {
	isBit0Set, _ := cpu.IsBitSet(*n8, 0)

	*n8 >>= 1
	if isBit0Set {
		*n8 |= 0b1000_0000
	}

	flags.ClearAll()
	if setFlagZ && *n8 == 0 {
		flags.SetZ()
	}
	if isBit0Set {
		flags.SetC()
	}
}

func LeftCircularAccumulator(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr) {
	leftCircularImpl(registerA, flags, false)

	instructionLengthInBytes := 1
	instruction := "RLCA"
	description := "RLCA: Rotate register A left."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LeftAccumulator(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr) {
	leftImpl(registerA, flags, false)

	instructionLengthInBytes := 1
	instruction := "RLA"
	description := "RLA: Rotate register A left, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightCircularAccumulator(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr) {
	rightCircularImpl(registerA, flags, false)

	instructionLengthInBytes := 1
	instruction := "RRCA"
	description := "RRCA: Rotate register A right."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightAccumulator(memory *cpu.Memory, programCounter uint16, registerA *uint8, flags cpu.FlagsPtr) {
	rightImpl(registerA, flags, false)

	instructionLengthInBytes := 1
	instruction := "RRA"
	description := "RRA: Rotate register A right, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LeftCircularRegister(memory *cpu.Memory, programCounter uint16, register *uint8, flags cpu.FlagsPtr, registerName string) {
	leftImpl(register, flags, true)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RLC %s", registerName)
	description := "RLC r8: Rotates register r8 left."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LeftCircularIndirectHL(memory *cpu.Memory, programCounter uint16, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	leftCircularImpl(&n8, flags, true)
	memory.Write(registerHL, n8)

	instructionLengthInBytes := 1
	instruction := "RLC [HC]"
	description := "RLC [HC]: Rotate the byte pointed to by HL left."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LeftRegister(memory *cpu.Memory, programCounter uint16, register *uint8, flags cpu.FlagsPtr, registerName string) {
	leftImpl(register, flags, true)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RL %s", registerName)
	description := "RL r8: Rotate bits in register r8 left, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LeftIndirectHL(memory *cpu.Memory, programCounter uint16, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	leftImpl(&n8, flags, true)
	memory.Write(registerHL, n8)

	instructionLengthInBytes := 1
	instruction := "RL [HL]"
	description := "RL [HL]: Rotate the byte pointed to by HL left, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightCircularRegister(memory *cpu.Memory, programCounter uint16, register *uint8, flags cpu.FlagsPtr, registerName string) {
	rightCircularImpl(register, flags, true)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RRC %s", registerName)
	description := "RRC r8: Rotates register r8 right."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightRegister(memory *cpu.Memory, programCounter uint16, register *uint8, flags cpu.FlagsPtr, registerName string) {
	rightImpl(register, flags, true)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RR %s", registerName)
	description := "RR r8: Rotate bits in register r8 right, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightIndirectHL(memory *cpu.Memory, programCounter uint16, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	rightImpl(&n8, flags, true)
	memory.Write(registerHL, n8)

	instructionLengthInBytes := 1
	instruction := "RR [HL]"
	description := "RR [HL]: Rotate the byte pointed to by HL right, through the carry flag."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func RightCircularIndirectHL(memory *cpu.Memory, programCounter uint16, registerHL uint16, flags cpu.FlagsPtr) {
	n8 := memory.Read(registerHL)
	rightCircularImpl(&n8, flags, true)
	memory.Write(registerHL, n8)

	instructionLengthInBytes := 1
	instruction := "RRC [HC]"
	description := "RRC [HC]: Rotate the byte pointed to by HL right."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}
