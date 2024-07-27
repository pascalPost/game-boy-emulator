package instructions

import (
	"fmt"
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions/utils"
	"log/slog"
)

func Nop(programCounter uint16) {
	slog.Debug("Instruction", "PC", cpu.FmtHex16(programCounter-1), "instruction", "NOP", "description", "No operation. This instruction doesn't do anything, but can be used to add a delay of one machine cycle and increment PC by one.")
}

func mostAndLeastSignificantByte(value uint16) (uint8, uint8) {
	// littleEndian
	leastSignificantByte := uint8(value & 0xff)
	mostSignificantByte := uint8(value >> 8)
	return mostSignificantByte, leastSignificantByte
}

func Jp(memory *cpu.Memory, programCounter *uint16) {
	pc := *programCounter
	instructionLengthInBytes := 3
	a16 := utils.ReadUnsigned16(memory, programCounter)
	*programCounter = a16
	cpu.Log(memory, pc, instructionLengthInBytes, fmt.Sprintf("JP 0x%04X", a16), "JP nn: Unconditional jump to the absolute address specified by the 16-bit immediate operand nn.")
}

func RelativeJump(memory *cpu.Memory, programCounter *uint16) {
	pc := *programCounter

	n8 := utils.ReadUnsigned8(memory, programCounter)
	e8 := int8(n8)
	newPc := int32(*programCounter) + int32(e8)
	if newPc < 0 {
		panic("negative program counter encountered.")
	}
	*programCounter = uint16(newPc)

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("JR %d", e8)
	description := "JR e: Unconditional jump to the relative address speciﬁed by the signed 8-bit operand e."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func RelativeJumpConditional(memory *cpu.Memory, programCounter *uint16, condition bool, conditionName string) {
	pc := *programCounter

	n8 := utils.ReadUnsigned8(memory, programCounter)
	e8 := int8(n8)
	if condition {
		newPc := int32(*programCounter) + int32(e8)
		if newPc < 0 {
			panic("negative program counter encountered.")
		}
		*programCounter = uint16(newPc)
	}

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("JR %s, %d (0x%04X)", conditionName, e8, *programCounter)
	description := "JR cc, e: Conditional jump to the relative address speciﬁed by the signed 8-bit operand e, depending on the condition cc."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func Load16BitToRegister(memory *cpu.Memory, programCounter *uint16, register *uint16, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 3
	n16 := utils.ReadUnsigned16(memory, programCounter)
	*register = n16
	cpu.Log(memory, pc, instructionLengthInBytes, fmt.Sprintf("LD %s, 0x%04X", registerName, n16), "LD rr, nn: Load to the 16-bit register rr, the immediate 16-bit data nn.")
}

func Load8BitToRegister(memory *cpu.Memory, programCounter *uint16, register *uint8, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 2
	n8 := utils.ReadUnsigned8(memory, programCounter)
	*register = n8
	instruction := fmt.Sprintf("LD %s, 0x%02X", registerName, n8)
	description := "LD r, n: Load to the 8-bit register r, the immediate 8-bit data n."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func Load8BitToRegisterFromRegister(memory *cpu.Memory, programCounter uint16, registerDest *uint8, registerSource uint8, registerDestName, registerSourceName string) {
	instructionLengthInBytes := 1
	*registerDest = registerSource
	instruction := fmt.Sprintf("LD %s, %s", registerDestName, registerSourceName)
	description := "LD r, r': Load to the 8-bit register r, data from the 8-bit register r'."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func Load8BitToRegisterFromAddressInRegister(memory *cpu.Memory, programCounter uint16, register *uint8, sourceRegisterValue uint16, registerName, sourceRegisterName string) {
	instructionLengthInBytes := 1
	*register = memory.Read(sourceRegisterValue)
	instruction := fmt.Sprintf("LD %s, [%s]", registerName, sourceRegisterName)
	description := "LD r, [r']: Load to the 8-bit register r, data from the absolute address speciﬁed by the 16-bit register r'."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func Load8BitToAddressInHLFromRegister(memory *cpu.Memory, programCounter uint16, register *uint8, registerHL uint16, registerName string) {
	instructionLengthInBytes := 1
	*register = memory.Read(registerHL)
	instruction := fmt.Sprintf("LD [HL], %s", registerName)
	description := "LD [HL], r: Load to the absolute address speciﬁed by the 16-bit register HL, data from the 8-bit register r."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func Load8BitToAddressInHL(memory *cpu.Memory, programCounter *uint16, registerHL uint16) {
	pc := *programCounter
	instructionLengthInBytes := 2
	n8 := utils.ReadUnsigned8(memory, programCounter)
	memory.Write(registerHL, n8)
	instruction := fmt.Sprintf("LD [HL], %02X", n8)
	description := "LD [HL], n: Load to the absolute address speciﬁed by the 16-bit register HL, the immediate data n."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func LoadFromAccumulatorDirect(memory *cpu.Memory, programCounter *uint16, registerA uint8) {
	pc := *programCounter
	a16 := utils.ReadUnsigned16(memory, programCounter)
	memory.Write(a16, registerA)

	instructionLengthInBytes := 3
	instruction := fmt.Sprintf("LD [%04X], A", a16)
	description := "LD [a16], A: Load to the absolute address speciﬁed by the 16-bit operand a16, data from the 8-bit A register."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func LoadFromAccumulatorDirectH(memory *cpu.Memory, programCounter *uint16, registerA uint8) {
	pc := *programCounter
	n8 := utils.ReadUnsigned8(memory, programCounter)
	mostSignificantByte := uint8(0xFF)
	a16 := utils.Unsigned16(n8, mostSignificantByte)
	memory.Write(a16, registerA)

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("LDH [%02X], A", n8)
	description := "LDH [n8], A: Load to the address speciﬁed by the 8-bit immediate data n, data from the 8-bit A register. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of n, so the possible range is 0xFF00-0xFFFF."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func LoadFromRegisterIndirect(memory *cpu.Memory, programCounter uint16, registerDest uint16, registerSource uint8, registerDestName, registerSourceName string) {
	a16 := registerDest
	memory.Write(a16, registerSource)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("LD [%s], %s", registerDestName, registerSourceName)
	description := "LD [r'], r: Load to the absolute address speciﬁed by the 16-bit register r', data from the 8-bit register r."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LoadAccumulatorDirectLeastSignificantByte(memory *cpu.Memory, programCounter *uint16, registerA *uint8) {
	pc := *programCounter
	n8 := utils.ReadUnsigned8(memory, programCounter)
	mostSignificantByte := uint8(0xFF)
	a16 := utils.Unsigned16(n8, mostSignificantByte)
	a := memory.Read(a16)
	*registerA = a

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("LDH A, [0x%02X]", n8)
	description := "LDH A, [n]: Load to the 8-bit A register, data from the address speciﬁed by the 8-bit immediate data n. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of n, so the possible range is 0xFF00-0xFFFF."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func LoadFromAccumulatorIndirectHC(memory *cpu.Memory, programCounter uint16, registerA uint8, registerC uint8) {
	mostSignificantByte := uint8(0xFF)
	a16 := utils.Unsigned16(registerC, mostSignificantByte)
	memory.Write(a16, registerA)

	instructionLengthInBytes := 1
	instruction := "LDH [C], A"
	description := "LDH [C], A: Load to the address speciﬁed by the 8-bit C register, data from the 8-bit A register. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of C, so the possible range is 0xFF00-0xFFFF."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LoadAccumulatorIndirectHLIncrement(memory *cpu.Memory, programCounter uint16, registerA *uint8, registerHL *uint16) {
	a8 := memory.Read(*registerHL)
	*registerA = a8
	*registerHL++

	instructionLengthInBytes := 1
	instruction := "LDI A, [HL]"
	description := "LDI A, [HL]: Load to the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL. The value of HL is incremented after the memory read."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LoadFromAccumulatorIndirectHLDecrement(memory *cpu.Memory, programCounter uint16, registerA uint8, registerHL *uint16) {
	n8 := registerA
	a16 := *registerHL
	memory.Write(a16, n8)
	*registerHL--

	instructionLengthInBytes := 1
	instruction := "LDD [HL], A"
	description := "LDD [HL], A: Load to the absolute address speciﬁed by the 16-bit register HL, data from the 8-bit A register. The value of HL is decremented after the memory write."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LoadFromAccumulatorIndirectHLIncrement(memory *cpu.Memory, programCounter uint16, registerA uint8, registerHL *uint16) {
	n8 := registerA
	a16 := *registerHL
	memory.Write(a16, n8)
	*registerHL++

	instructionLengthInBytes := 1
	instruction := "LDI [HL], A"
	description := "LDI [HL], A: Load to the absolute address speciﬁed by the 16-bit register HL, data from the 8-bit A register. The value of HL is incremented after the memory write."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func LoadAccumulatorIndirectHLDecrement(memory *cpu.Memory, programCounter uint16, registerA *uint8, registerHL *uint16) {
	a8 := memory.Read(*registerHL)
	*registerA = a8
	*registerHL--

	instructionLengthInBytes := 1
	instruction := "LDD A, [HL]"
	description := "LDD A, [HL]: Load to the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL. The value of HL is incremented after the memory read."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func push(memory *cpu.Memory, stackPointer *uint16, address uint16) {
	*stackPointer--
	msb, lsb := mostAndLeastSignificantByte(address)
	memory.Write(*stackPointer, msb)
	*stackPointer--
	memory.Write(*stackPointer, lsb)
	slog.Debug("Push to stack", "address", cpu.FmtHex16(address), "new stack pointer address", cpu.FmtHex16(*stackPointer))
}

func Call(memory *cpu.Memory, programCounter *uint16, stackPointer *uint16) {
	pc := *programCounter
	a16 := utils.ReadUnsigned16(memory, programCounter)
	addressOfNextInstruction := *programCounter
	push(memory, stackPointer, addressOfNextInstruction)
	*programCounter = a16

	instructionLengthInBytes := 3
	instruction := fmt.Sprintf("CALL 0x%04X", a16)
	description := "CALL nn: Unconditional function call to the absolute address specified by the 16-bit operand nn."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func PushToStack(memory *cpu.Memory, programCounter uint16, stackPointer *uint16, register uint16, registerName string) {
	push(memory, stackPointer, register)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("PUSH %s", registerName)
	description := "PUSH rr: Push to the stack memory, data from the 16-bit register rr."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func PopFromStack(memory *cpu.Memory, programCounter uint16, stackPointer *uint16, register *uint16, registerName string) {
	leastSignificantByte := memory.Read(*stackPointer)
	*stackPointer++
	mostSignificantByte := memory.Read(*stackPointer)
	*stackPointer++
	*register = utils.Unsigned16(leastSignificantByte, mostSignificantByte)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("POP %s", registerName)
	description := "POP rr: Pops to the 16-bit register rr, data from the stack memory. This instruction does not do calculations that aﬀect ﬂags, but POP AF completely replaces the F register value, so all ﬂags are changed based on the 8-bit data that is read from memory."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func returnImpl(memory *cpu.Memory, programCounter *uint16, stackPointer *uint16) {
	leastSignificantByte := memory.Read(*stackPointer)
	*stackPointer++
	mostSignificantByte := memory.Read(*stackPointer)
	*stackPointer++
	a16 := utils.Unsigned16(leastSignificantByte, mostSignificantByte)
	*programCounter = a16
}

func ReturnFromFunction(memory *cpu.Memory, programCounter *uint16, stackPointer *uint16) {
	pc := *programCounter

	returnImpl(memory, programCounter, stackPointer)

	instructionLengthInBytes := 1
	instruction := "RET"
	description := "RET: Unconditional return from a function."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func ReturnFromFunctionConditional(memory *cpu.Memory, programCounter *uint16, stackPointer *uint16, condition bool, conditionName string) {
	pc := *programCounter

	if condition {
		returnImpl(memory, programCounter, stackPointer)
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RET %s", conditionName)
	description := "RET cc: Conditional return from a function, depending on the condition cc."
	cpu.Log(memory, pc, instructionLengthInBytes, instruction, description)
}

func DisableInterrupts(memory *cpu.Memory, programCounter uint16, ime *bool) {
	*ime = false

	instructionLengthInBytes := 1
	instruction := "DI"
	description := "Disables interrupt handling by setting IME=0 and cancelling any scheduled eﬀects of the EI instruction if any."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func TestBitRegister(memory *cpu.Memory, programCounter uint16, flags cpu.FlagsPtr, bitNumber uint8, register uint8, registerName string) {
	bitSet, err := cpu.IsBitSet(register, bitNumber)
	if err != nil {
		panic(err)
	}

	if bitSet {
		flags.ClearZ()
	} else {
		flags.SetZ()
	}
	flags.ClearN()
	flags.SetH()

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("BIT %d, %s", bitNumber, registerName)
	description := "BIT b, r: Test bit b n register r."
	cpu.Log(memory, programCounter, instructionLengthInBytes, instruction, description)
}
