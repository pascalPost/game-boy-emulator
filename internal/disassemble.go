package internal

import (
	"fmt"
	"log"
	"slices"
)

func isPrefixed(b byte) bool {
	return b == 0xCB
}

func toLittleEndian(b []byte) []byte {
	address := make([]byte, len(b))
	copy(address, b)
	slices.Reverse(address)
	return address
}

func littleEndian16BitAddressOrData(data []byte, programCounter uint16, operand *Operand) (uint16, string) {
	if operand.Bytes != 2 {
		log.Fatal("unexpected number of bytes.")
	}
	value := toLittleEndian(data[programCounter : programCounter+2])
	str := fmt.Sprintf("0x%X", value)
	programCounter += 2
	return programCounter, str
}

func unsigned8BitData(data []byte, programCounter uint16, operand *Operand) (uint16, string) {
	if operand.Bytes != 1 {
		log.Fatal("unexpected number of bytes.")
	}
	str := fmt.Sprintf("0x%.2X", data[programCounter])
	programCounter++
	return programCounter, str
}

func signed8BitData(data []byte, programCounter uint16, operand *Operand) (uint16, string) {
	if operand.Bytes != 1 {
		log.Fatal("unexpected number of bytes.")
	}
	str := fmt.Sprintf("%d (0x%.2X)", int8(data[programCounter]), data[programCounter])
	programCounter++
	return programCounter, str
}

func handleOperand(data []byte, programCounter uint16, operand *Operand) (uint16, string) {
	operandStr := ""
	switch operand.Name {
	case "0":
		operandStr += "0"
	case "1":
		operandStr += "1"
	case "2":
		operandStr += "2"
	case "3":
		operandStr += "3"
	case "4":
		operandStr += "4"
	case "5":
		operandStr += "5"
	case "6":
		operandStr += "6"
	case "7":
		operandStr += "7"
	case "AF":
		// 16-bit register (Accumulator & Flags)
		operandStr += "AF"
	case "A":
		// 8-bit Hi part of AF
		operandStr += "A"
	case "BC":
		// 16-bit register
		operandStr += "BC"
	case "B":
		// 8-bit Hi part of BC
		operandStr += "B"
	case "C":
		// 8-bit Lo part of BC
		// or
		// Condition code: Execute if C is set.
		operandStr += "C"
	case "DE":
		// 16-bit register
		operandStr += "DE"
	case "D":
		// 8-bit Hi part of DE
		operandStr += "D"
	case "E":
		// 8-bit Lo part of DE
		operandStr += "E"
	case "HL":
		// 8-bit Hi part of HL
		operandStr += "HL"
	case "H":
		// 8-bit Hi part of HL
		operandStr += "H"
	case "L":
		// 8-bit Lo part of HL
		operandStr += "L"
	case "SP":
		// 16-bit register (Stack Pointer)
		operandStr += "SP"
	case "n8":
		// immediate 8-bit data
		return unsigned8BitData(data, programCounter, operand)
	case "n16":
		// immediate little-endian 16-bit data
		return littleEndian16BitAddressOrData(data, programCounter, operand)
	case "a8":
		// means 8-bit unsigned data, which is added to $FF00 in certain instructions to create a 16-bit address in HRAM (High RAM)
		return unsigned8BitData(data, programCounter, operand)
	case "a16":
		// little-endian 16-bit address
		return littleEndian16BitAddressOrData(data, programCounter, operand)
	case "e8":
		// 8-bit signed data (offset)
		return signed8BitData(data, programCounter, operand)
	case "Z":
		// Condition code: Execute if Z is set.
		operandStr += "Z"
	case "NZ":
		// Condition code: Execute if Z is not set.
		operandStr += "NZ"
	case "NC":
		// Condition code: Execute if C is not set.
		operandStr += "NC"
	case "$00":
		operandStr += "0x00(H)"
	case "$08":
		operandStr += "0x08(H)"
	case "$10":
		operandStr += "0x10(H)"
	case "$18":
		operandStr += "0x18(H)"
	case "$20":
		operandStr += "0x20(H)"
	case "$28":
		operandStr += "0x28(H)"
	case "$30":
		operandStr += "0x30(H)"
	case "$38":
		operandStr += "0x38(H)"
	default:
		log.Panicf("unknown operand name: {%s}", operand.Name)
	}
	return programCounter, operandStr
}

func readOperands(data []byte, programCounter uint16, operands []Operand) (uint16, string) {
	operandStr := ""
	for i, operand := range operands {
		newProgramCounter, str := handleOperand(data, programCounter, &operand)
		programCounter = newProgramCounter
		if !operand.Immediate {
			str = fmt.Sprintf("[%s]", str)
		}
		if i < len(operands)-1 {
			str = fmt.Sprintf("%s, ", str)
		}
		operandStr += str
	}
	return programCounter, fmt.Sprintf(" %s", operandStr)
}

type Instruction struct {
	Line         string
	AddressStart uint16
	AddressEnd   uint16
}

func parseOpcode(data []byte, programCounter uint16, list *OpcodeList) (uint16, Opcode) {
	prefixed := isPrefixed(data[programCounter])

	var opcode Opcode
	if !prefixed {
		opcode = list.UnPrefixed[ByteKey{data[programCounter]}]
	} else {
		programCounter++
		opcode = list.CbPrefixed[ByteKey{data[programCounter]}]
	}
	programCounter++
	return programCounter, opcode
}

func Disassemble(data []byte, programCounter uint16, list *OpcodeList) []Instruction {
	instructions := make([]Instruction, 0, len(data))

	for int(programCounter) < len(data) {
		i := Instruction{}
		i.AddressStart = programCounter

		newProgramCounter, opcode := parseOpcode(data, programCounter, list)
		programCounter = newProgramCounter
		i.Line += opcode.Mnemonic

		newProgramCounter, str := readOperands(data, programCounter, opcode.Operands)
		programCounter = newProgramCounter
		i.Line += str

		i.AddressEnd = programCounter
		instructions = append(instructions, i)
	}

	return instructions
}
