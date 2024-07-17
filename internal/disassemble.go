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
	str := fmt.Sprintf(" 0x%X", value)
	programCounter += 2
	return programCounter, str
}

func immediate8BitData(data []byte, programCounter uint16, operand *Operand) (uint16, string) {
	if operand.Bytes != 1 {
		log.Fatal("unexpected number of bytes.")
	}
	str := fmt.Sprintf(" 0x%.2X", data[programCounter])
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
	case "AF":
		// 16-bit register (Accumulator & Flags)
		operandStr += " AF"
	case "A":
		// 8-bit Hi part of AF
		operandStr += " A"
	case "BC":
		// 8-bit Hi part of BC
		operandStr += " BC"
	case "B":
		// 8-bit Hi part of BC
		operandStr += " B"
	case "C":
		// 8-bit Lo part of BC
		// or
		// Condition code: Execute if C is set.
		operandStr += " C"
	case "D":
		// 8-bit Hi part of DE
		operandStr += " D"
	case "E":
		// 8-bit Lo part of DE
		operandStr += " E"
	case "HL":
		// 8-bit Hi part of HL
		operandStr += " HL"
	case "H":
		// 8-bit Hi part of HL
		operandStr += " H"
	case "L":
		// 8-bit Lo part of HL
		operandStr += " L"
	case "SP":
		// 16-bit register (Stack Pointer)
		operandStr += " SP"
	case "n8":
		// immediate 8-bit data
		return immediate8BitData(data, programCounter, operand)
	case "n16":
		// immediate little-endian 16-bit data
		return littleEndian16BitAddressOrData(data, programCounter, operand)
	case "a16":
		// little-endian 16-bit address
		return littleEndian16BitAddressOrData(data, programCounter, operand)
	case "e8":
		// 8-bit signed data (offset)
		return signed8BitData(data, programCounter, operand)
	case "Z":
		// Condition code: Execute if Z is set.
		operandStr += " Z"
	case "NZ":
		// Condition code: Execute if Z is not set.
		operandStr += " NZ"
	case "NC":
		// Condition code: Execute if C is not set.
		operandStr += " NC"
	default:
		log.Panicf("unknown operand name: {%s}", operand.Name)
	}
	return programCounter, operandStr
}

func readOperands(data []byte, programCounter uint16, operands []Operand) (uint16, string) {
	operandStr := ""
	for _, operand := range operands {
		newProgramCounter, str := handleOperand(data, programCounter, &operand)
		programCounter = newProgramCounter
		operandStr += str
	}
	return programCounter, operandStr
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
		opcode = list.UnPrefixed[ByteKey{data[programCounter]}]
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
