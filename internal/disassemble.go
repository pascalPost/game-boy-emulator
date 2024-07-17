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

func littleEndian16BitAddressOrData(data []byte, operand *Operand) ([]byte, string) {
	if operand.Bytes != 2 {
		log.Fatal("unexpected number of bytes.")
	}
	value := toLittleEndian(data[0:2])
	str := fmt.Sprintf(" 0x%X", value)
	remainingData := data[2:]
	return remainingData, str
}

func immediate8BitData(data []byte, operand *Operand) ([]byte, string) {
	if operand.Bytes != 1 {
		log.Fatal("unexpected number of bytes.")
	}
	str := fmt.Sprintf(" 0x%.2X", data[0])
	remainingData := data[1:]
	return remainingData, str
}

func handleOperand(data []byte, operand *Operand) ([]byte, string) {
	switch operand.Name {
	case "AF":
		// 16-bit register (Accumulator & Flags)
		return data, " AF"
	case "A":
		// 8-bit Hi part of AF
		return data, " A"
	case "B":
		// 8-bit Hi part of BC
		return data, " B"
	case "C":
		// 8-bit Lo part of BC
		return data, " C"
	case "D":
		// 8-bit Hi part of DE
		return data, " D"
	case "E":
		// 8-bit Lo part of DE
		return data, " E"
	case "H":
		// 8-bit Hi part of HL
		return data, " H"
	case "L":
		// 8-bit Lo part of HL
		return data, " L"
	case "SP":
		// 16-bit register (Stack Pointer)
		return data, " SP"
	case "n8":
		// immediate 8-bit data
		return immediate8BitData(data, operand)
	case "n16":
		// immediate little-endian 16-bit data
		return littleEndian16BitAddressOrData(data, operand)
	case "a16":
		// little-endian 16-bit address
		return littleEndian16BitAddressOrData(data, operand)
	default:
		log.Panicf("unknown operand name: {%s}", operand.Name)
	}
	return nil, ""
}

func readOperands(data []byte, operands []Operand) ([]byte, string) {
	operandStr := ""
	for _, operand := range operands {
		newDataStack, str := handleOperand(data, &operand)
		data = newDataStack
		operandStr += str
	}
	return data, operandStr
}

type Instruction struct {
	Line    string
	Address byte
	Data    []byte
}

func Disassemble(dataStack []byte, list *OpcodeList) []Instruction {
	instructions := make([]Instruction, 0, len(dataStack))

	for len(dataStack) > 0 {
		i := Instruction{}

		prefixed := isPrefixed(dataStack[0])

		var opcode Opcode
		if !prefixed {
			opcode = list.UnPrefixed[ByteKey{dataStack[0]}]
		} else {
			dataStack = dataStack[1:]
			opcode = list.UnPrefixed[ByteKey{dataStack[0]}]
		}
		dataStack = dataStack[1:]

		i.Line += opcode.Mnemonic

		data, str := readOperands(dataStack, opcode.Operands)
		dataStack = data
		i.Line += str

		instructions = append(instructions, i)
	}

	return instructions
}
