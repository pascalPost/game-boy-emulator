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

func loadTo16BitRegister(memory *memory, programCounter *uint16, register *uint16, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 3
	n16 := readUnsigned16(memory, programCounter)
	*register = n16
	logInstruction(memory, pc, instructionLengthInBytes, fmt.Sprintf("LD %s 0x%04X", registerName, n16), "LD rr, nn: Load to the 16-bit register rr, the immediate 16-bit data nn.")
}

func loadTo8BitRegister(memory *memory, programCounter *uint16, register *uint8, registerName string) {
	pc := *programCounter
	instructionLengthInBytes := 2
	n8 := readUnsigned8(memory, programCounter)
	*register = n8
	logInstruction(memory, pc, instructionLengthInBytes, fmt.Sprintf("LD %s 0x%02X", registerName, n8), "LD r, n: Load to the 8-bit register r, the immediate 8-bit data n.")
}

func push(memory *memory, stackPointer *uint16, address uint16) {
	*stackPointer--
	msb, lsb := mostAndLeastSignificantByte(address)
	memory.write(*stackPointer, msb)
	*stackPointer--
	memory.write(*stackPointer, lsb)
	slog.Debug("Push to stack", "address", fmtHex16(address), "new stack pointer address", fmtHex16(*stackPointer))
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
