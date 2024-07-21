package internal

import (
	"fmt"
	"log/slog"
)

func logInstruction(memory *memory, programCounter uint16, instructionLengthInBytes int, instruction, description string) {
	pcBegin := programCounter - 1
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

func loadFromAccumulator(memory *memory, programCounter *uint16, registerA uint8) {
	pc := *programCounter
	instructionLengthInBytes := 3
	a16 := readUnsigned16(memory, programCounter)
	memory.write(a16, registerA)
	instruction := fmt.Sprintf("LD %04X A", a16)
	description := "LD a16, A: Load to the absolute address speciﬁed by the 16-bit operand a16, data from the 8-bit A register."
	logInstruction(memory, pc, instructionLengthInBytes, instruction, description)
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

func returnFromFunctionConditional(memory *memory, programCounter *uint16, stackPointer *uint16,
	condition bool,
	conditionName string) {
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

	flags.clear()
	if result == 0 {
		flags.setZ()
	}

	instructionLengthInBytes := 1
	instruction := fmt.Sprintf("OR A %s", registerName)
	description := "Performs a bitwise OR operation between the 8-bit A register and the 8-bit register r, and stores the result back into the A register."
	logInstruction(memory, programCounter, instructionLengthInBytes, instruction, description)
}
