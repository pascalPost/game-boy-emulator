package internal

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log/slog"
)

func logInstruction(memory *memory, programCounter uint16, instructionLengthInBytes int, instruction, description string) {
	pcBegin := programCounter - 1

	// handle prefixed instruction
	if pcBegin > 0 && memory.data[pcBegin-1] == uint8(0xCB) {
		pcBegin--
	}

	pcEnd := pcBegin + uint16(instructionLengthInBytes)
	slog.Debug("Instruction", "PC", fmtHex16(pcBegin), "mem", fmt.Sprintf("0x% 2X", memory.data[pcBegin:pcEnd]), "instruction", instruction, "description", description)
}

func nop(programCounter uint16) {
	slog.Debug("Instruction", "PC", fmtHex16(programCounter-1), "instruction", "NOP", "description", "No operation. This instruction doesn't do anything, but can be used to add a delay of one machine cycle and increment PC by one.")
}

func unsigned16(leastSignificantByte uint8, mostSignificantByte uint8) uint16 {
	// littleEndian
	return uint16(leastSignificantByte) | uint16(mostSignificantByte)<<8
}

func mostAndLeastSignificantByte(value uint16) (uint8, uint8) {
	// littleEndian
	leastSignificantByte := uint8(value & 0xff)
	mostSignificantByte := uint8(value >> 8)
	return mostSignificantByte, leastSignificantByte
}

func fmtHex16(value uint16) string {
	return fmt.Sprintf("0x%04X", value)
}

func fmtHex8(value uint8) string {
	return fmt.Sprintf("0x%02X", value)
}

func readUnsigned8(memory *memory, programCounter *uint16) uint8 {
	v := memory.read(*programCounter)
	*programCounter++
	return v
}

func readUnsigned16(memory *memory, programCounter *uint16) uint16 {
	nnLSB := readUnsigned8(memory, programCounter)
	nnMSB := readUnsigned8(memory, programCounter)
	nn := unsigned16(nnLSB, nnMSB)
	return nn
}

func jp(memory *memory, programCounter *uint16) {
	pc := *programCounter
	instructionLengthInBytes := 3
	a16 := readUnsigned16(memory, programCounter)
	*programCounter = a16
	logInstruction(memory, pc, instructionLengthInBytes, fmt.Sprintf("JP 0x%04X", a16), "JP nn: Unconditional jump to the absolute address specified by the 16-bit immediate operand nn.")
}

func relativeJump(memory *memory, programCounter *uint16) {
	pc := *programCounter

	n8 := readUnsigned8(memory, programCounter)
	e8 := int8(n8)
	newPc := int32(*programCounter) + int32(e8)
	if newPc < 0 {
		panic("negative program counter encountered.")
	}
	*programCounter = uint16(newPc)

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("JR %d", e8)
	description := "JR e: Unconditional jump to the relative address speciﬁed by the signed 8-bit operand e."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func relativeJumpConditional(memory *memory, programCounter *uint16, condition bool, conditionName string) {
	pc := *programCounter

	n8 := readUnsigned8(memory, programCounter)
	e8 := int8(n8)
	if condition {
		newPc := int32(*programCounter) + int32(e8)
		if newPc < 0 {
			panic("negative program counter encountered.")
		}
		*programCounter = uint16(newPc)
	}

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("JR %s, %d", conditionName, e8)
	description := "JR cc, e: Conditional jump to the relative address speciﬁed by the signed 8-bit operand e, depending on the condition cc."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func load16BitToRegister(memory *memory, programCounter *uint16, register *uint16, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 3
	n16 := readUnsigned16(memory, programCounter)
	*register = n16
	logInstruction(memory, pc, instructionLengthInBytes, fmt.Sprintf("LD %s, 0x%04X", registerName, n16), "LD rr, nn: Load to the 16-bit register rr, the immediate 16-bit data nn.")
}

func load8BitToRegister(memory *memory, programCounter *uint16, register *uint8, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 2
	n8 := readUnsigned8(memory, programCounter)
	*register = n8
	instruction := fmt.Sprintf("LD %s, 0x%02X", registerName, n8)
	description := "LD r, n: Load to the 8-bit register r, the immediate 8-bit data n."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func load8BitToRegisterFromRegister(memory *memory, programCounter uint16, registerDest *uint8, registerSource uint8, registerDestName, registerSourceName string) {
	instructionLengthInBytes := 1
	*registerDest = registerSource
	instruction := fmt.Sprintf("LD %s, %s", registerDestName, registerSourceName)
	description := "LD r, r': Load to the 8-bit register r, data from the 8-bit register r'."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func load8BitToRegisterFromAddressInRegister(memory *memory, programCounter uint16, register *uint8, sourceRegisterValue uint16, registerName, sourceRegisterName string) {
	instructionLengthInBytes := 1
	*register = memory.read(sourceRegisterValue)
	instruction := fmt.Sprintf("LD %s, [%s]", registerName, sourceRegisterName)
	description := "LD r, [r']: Load to the 8-bit register r, data from the absolute address speciﬁed by the 16-bit register r'."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func load8BitToAddressInHLFromRegister(memory *memory, programCounter uint16, register *uint8, registerHL uint16, registerName string) {
	instructionLengthInBytes := 1
	*register = memory.read(registerHL)
	instruction := fmt.Sprintf("LD [HL], %s", registerName)
	description := "LD [HL], r: Load to the absolute address speciﬁed by the 16-bit register HL, data from the 8-bit register r."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func load8BitToAddressInHL(memory *memory, programCounter *uint16, registerHL uint16) {
	pc := *programCounter
	instructionLengthInBytes := 2
	n8 := readUnsigned8(memory, programCounter)
	memory.write(registerHL, n8)
	instruction := fmt.Sprintf("LD [HL], %02X", n8)
	description := "LD [HL], n: Load to the absolute address speciﬁed by the 16-bit register HL, the immediate data n."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func loadFromAccumulatorDirect(memory *memory, programCounter *uint16, registerA uint8) {
	pc := *programCounter
	a16 := readUnsigned16(memory, programCounter)
	memory.write(a16, registerA)

	instructionLengthInBytes := 3
	instruction := fmt.Sprintf("LD [%04X], A", a16)
	description := "LD [a16], A: Load to the absolute address speciﬁed by the 16-bit operand a16, data from the 8-bit A register."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func loadFromAccumulatorDirectH(memory *memory, programCounter *uint16, registerA uint8) {
	pc := *programCounter
	n8 := readUnsigned8(memory, programCounter)
	mostSignificantByte := uint8(0xFF)
	a16 := unsigned16(n8, mostSignificantByte)
	memory.write(a16, registerA)

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("LDH [%02X], A", n8)
	description := "LDH [n8], A: Load to the address speciﬁed by the 8-bit immediate data n, data from the 8-bit A register. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of n, so the possible range is 0xFF00-0xFFFF."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func loadFromRegisterIndirect(memory *memory, programCounter uint16, registerDest uint16, registerSource uint8, registerDestName, registerSourceName string) {
	a16 := registerDest
	memory.write(a16, registerSource)

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("LD [%s], %s", registerDestName, registerSourceName)
	description := "LD [r'], r: Load to the absolute address speciﬁed by the 16-bit register r', data from the 8-bit register r."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func loadAccumulatorDirectLeastSignificantByte(memory *memory, programCounter *uint16, registerA *uint8) {
	pc := *programCounter
	n8 := readUnsigned8(memory, programCounter)
	mostSignificantByte := uint8(0xFF)
	a16 := unsigned16(n8, mostSignificantByte)
	a := memory.read(a16)
	*registerA = a

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("LDH A, [0x%02X]", n8)
	description := "LDH A, [n]: Load to the 8-bit A register, data from the address speciﬁed by the 8-bit immediate data n. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of n, so the possible range is 0xFF00-0xFFFF."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func loadFromAccumulatorIndirectHC(memory *memory, programCounter uint16, registerA uint8, registerC uint8) {
	mostSignificantByte := uint8(0xFF)
	a16 := unsigned16(registerC, mostSignificantByte)
	memory.write(a16, registerA)

	instructionLengthInBytes := 1
	instruction := "LDH [C], A"
	description := "LDH [C], A: Load to the address speciﬁed by the 8-bit C register, data from the 8-bit A register. The full 16-bit absolute address is obtained by setting the most signiﬁcant byte to 0xFF and the least signiﬁcant byte to the value of C, so the possible range is 0xFF00-0xFFFF."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func loadAccumulatorIndirectHLIncrement(memory *memory, programCounter uint16, registerA *uint8, registerHL *uint16) {
	a8 := memory.read(*registerHL)
	*registerA = a8
	*registerHL++

	instructionLengthInBytes := 1
	instruction := "LDI A, [HL]"
	description := "LDI A, [HL]: Load to the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL. The value of HL is incremented after the memory read."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func loadFromAccumulatorIndirectHLDecrement(memory *memory, programCounter uint16, registerA uint8, registerHL *uint16) {
	n8 := registerA
	a16 := *registerHL
	memory.write(a16, n8)
	*registerHL--

	instructionLengthInBytes := 1
	instruction := "LDD [HL], A"
	description := "LDD [HL], A: Load to the absolute address speciﬁed by the 16-bit register HL, data from the 8-bit A register. The value of HL is decremented after the memory write."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func loadAccumulatorIndirectHLDecrement(memory *memory, programCounter uint16, registerA *uint8, registerHL *uint16) {
	a8 := memory.read(*registerHL)
	*registerA = a8
	*registerHL--

	instructionLengthInBytes := 1
	instruction := "LDD A, [HL]"
	description := "LDD A, [HL]: Load to the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL. The value of HL is incremented after the memory read."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func call(memory *memory, programCounter *uint16, stackPointer *uint16) {
	pc := *programCounter
	instructionLengthInBytes := 3
	a16 := readUnsigned16(memory, programCounter)
	addressOfNextInstruction := *programCounter
	push(memory, stackPointer, addressOfNextInstruction)
	*programCounter = a16
	logInstruction(memory, pc, instructionLengthInBytes, fmt.Sprintf("CALL 0x%04X", a16), "CALL nn: Unconditional function call to the absolute address specified by the 16-bit operand nn.")
}

func push(memory *memory, stackPointer *uint16, address uint16) {
	*stackPointer--
	msb, lsb := mostAndLeastSignificantByte(address)
	memory.write(*stackPointer, msb)
	*stackPointer--
	memory.write(*stackPointer, lsb)
	slog.Debug("Push to stack", "address", fmtHex16(address), "new stack pointer address", fmtHex16(*stackPointer))
}

func returnImpl(memory *memory, programCounter *uint16, stackPointer *uint16) {
	leastSignificantByte := memory.read(*stackPointer)
	*stackPointer++
	mostSignificantByte := memory.read(*stackPointer)
	*stackPointer++
	a16 := unsigned16(leastSignificantByte, mostSignificantByte)
	*programCounter = a16
}

func returnFromFunction(memory *memory, programCounter *uint16, stackPointer *uint16) {
	pc := *programCounter

	returnImpl(memory, programCounter, stackPointer)

	instructionLengthInBytes := 1
	instruction := "RET"
	description := "RET: Unconditional return from a function."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func returnFromFunctionConditional(memory *memory, programCounter *uint16, stackPointer *uint16, condition bool, conditionName string) {
	pc := *programCounter

	if condition {
		returnImpl(memory, programCounter, stackPointer)
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("RET %s", conditionName)
	description := "RET cc: Conditional return from a function, depending on the condition cc."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func bitwiseOrRegister(memory *memory, programCounter uint16, registerA *uint8, flags flagsPtr, register uint8, registerName string) {
	result := *registerA | register
	*registerA = result

	flags.clearAll()
	if result == 0 {
		flags.setZ()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("OR A, %s", registerName)
	description := "Performs a bitwise OR operation between the 8-bit A register and the 8-bit register r, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func bitwiseXorRegister(memory *memory, programCounter uint16, registerA *uint8, flags flagsPtr, register uint8, registerName string) {
	result := *registerA ^ register
	*registerA = result

	flags.clearAll()
	if result == 0 {
		flags.setZ()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("XOR A, %s", registerName)
	description := "XOR r: Performs a bitwise XOR operation between the 8-bit A register and the 8-bit register r, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func halfCarryAdd(a, b uint8) bool {
	return (a&0b1111)+(b&0b1111) > 0b1111
}

func carryAdd(a, b uint8) bool {
	return uint16(a)+uint16(b) > 0b1111_1111
}

func addImpl(registerA *uint8, n8 uint8, flags flagsPtr) {
	a := *registerA
	result := a + n8
	*registerA = result

	flags.clearAll()

	if result == 0 {
		flags.setZ()
	}
	flags.setN()
	if halfCarryAdd(a, n8) {
		flags.setH()
	}
	if carryAdd(a, n8) {
		flags.setC()
	}
}

func addRegister(memory *memory, programCounter uint16, registerA *uint8, register uint8, flags flagsPtr, registerName string) {
	addImpl(registerA, register, flags)
	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("ADD %s", registerName)
	description := "Adds to the 8-bit A register, the 8-bit register r, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func addIndirectHL(memory *memory, programCounter uint16, registerA *uint8, registerHL uint16, flags flagsPtr) {
	n8 := memory.read(registerHL)
	addImpl(registerA, n8, flags)

	instructionLengthInBytes := 1
	instruction := "ADD [HL]"
	description := "Adds from the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func addImmediate(memory *memory, programCounter *uint16, registerA *uint8, flags flagsPtr) {
	pc := *programCounter
	n8 := readUnsigned8(memory, programCounter)
	addImpl(registerA, n8, flags)
	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("ADD 0x%02X", n8)
	description := "Adds from the 8-bit A register, the immediate data n, and stores the result back into the A register."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func halfCarrySub(a, b uint8) bool {
	return (a & 0b1111) < (b & 0b1111)
}

func carrySub(a, b uint8) bool {
	return a < b
}

func subtractImpl(registerA uint8, n8 uint8, flags flagsPtr) uint8 {
	result := registerA - n8

	flags.clearAll()

	if result == 0 {
		flags.setZ()
	}
	flags.setN()
	if halfCarrySub(registerA, n8) {
		flags.setH()
	}
	if carrySub(registerA, n8) {
		flags.setC()
	}

	return result
}

func subtractRegister(memory *memory, programCounter uint16, registerA *uint8, register uint8, flags flagsPtr, registerName string) {
	res := subtractImpl(*registerA, register, flags)
	*registerA = res
	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("SUB %s", registerName)
	description := "Subtracts from the 8-bit A register, the 8-bit register r, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func subtractIndirectHL(memory *memory, programCounter uint16, registerA *uint8, registerHL uint16, flags flagsPtr) {
	n8 := memory.read(registerHL)
	res := subtractImpl(*registerA, n8, flags)
	*registerA = res

	instructionLengthInBytes := 1
	instruction := "SUB [HL]"
	description := "Subtracts from the 8-bit A register, data from the absolute address speciﬁed by the 16-bit register HL, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func subtractImmediate(memory *memory, programCounter *uint16, registerA *uint8, flags flagsPtr) {
	pc := *programCounter

	n8 := readUnsigned8(memory, programCounter)
	res := subtractImpl(*registerA, n8, flags)
	*registerA = res

	instructionLengthInBytes := 2
	instruction := fmt.Sprintf("SUB 0x%02X", n8)
	description := "Subtracts from the 8-bit A register, the immediate data n, and stores the result back into the A register."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func incrementRegister(memory *memory, programCounter uint16, register *uint8, flags flagsPtr, registerName string) {
	old := *register
	*register = old + 1

	if *register == 0 {
		flags.setZ()
	} else {
		flags.clearZ()
	}
	flags.clearN()
	if halfCarryAdd(old, 1) {
		flags.setH()
	} else {
		flags.clearH()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("INC %s", registerName)
	description := "INC rr: Increments data in the 16-bit register rr."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func increment16BitRegister(memory *memory, programCounter uint16, register *uint16, registerName string) {
	*register++

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("INC %s", registerName)
	description := "INC rr: Increments data in the 16-bit register rr."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func decrement16BitRegister(memory *memory, programCounter uint16, register *uint16, registerName string) {
	*register--

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("DEC %s", registerName)
	description := "DEC rr: Decrements data in the 16-bit register rr."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func disableInterrupts(memory *memory, programCounter uint16, ime *bool) {
	*ime = false

	instructionLengthInBytes := 1
	instruction := "DI"
	description := "Disables interrupt handling by setting IME=0 and cancelling any scheduled eﬀects of the EI instruction if any."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}

func compareImmediate(memory *memory, programCounter *uint16, registerA uint8, flags flagsPtr) {
	pc := *programCounter

	n8 := readUnsigned8(memory, programCounter)
	subtractImpl(registerA, n8, flags)

	instructionLengthInBytes := 2
	instruction := spew.Sprintf("CP 0x%02X", n8)
	description := "Subtracts from the 8-bit A register, the immediate data n, and updates ﬂags based on the result. This instruction is basically identical to SUB n, but does not update the A register."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
}

func testBitRegister(memory *memory, programCounter uint16, flags flagsPtr, bitNumber uint8, register uint8, registerName string) {
	bitSet, err := isBitSet(register, bitNumber)
	if err != nil {
		panic(err)
	}

	if bitSet {
		flags.clearZ()
	} else {
		flags.setZ()
	}
	flags.clearN()
	flags.setH()

	instructionLengthInBytes := 1
	instruction := spew.Sprintf("BIT %d, %s", bitNumber, registerName)
	description := "BIT b, r: Test bit b n register r."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}
