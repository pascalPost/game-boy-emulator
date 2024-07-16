package internal

import (
	"fmt"
	"slices"
)

func isPrefixed(b byte) bool {
	return b == 0xCB
}

func toAddress(b []byte) []byte {
	address := make([]byte, len(b))
	copy(address, b)
	slices.Reverse(address)
	return address
}

func readOperands(data []byte, operands []Operand) []byte {
	for _, o := range operands {
		address := toAddress(data[0:o.Bytes])
		fmt.Printf(" 0x%X", address)
		data = data[o.Bytes:]
	}

	return data
}

type Instruction struct {
	Line    string
	Address byte
	Data    []byte
}

func Disassemble(dataStack []byte, list *OpcodeList) []Instruction {
	instructions := make([]Instruction, 0, len(dataStack))

	for len(dataStack) > 0 {
		prefixed := isPrefixed(dataStack[0])

		var opcode Opcode
		if !prefixed {
			opcode = list.UnPrefixed[ByteKey{dataStack[0]}]
		} else {
			//	i += 1
			dataStack = dataStack[1:]
			opcode = list.UnPrefixed[ByteKey{dataStack[0]}]
		}
		dataStack = dataStack[1:]

		fmt.Printf("%s", opcode.Mnemonic)

		dataStack = readOperands(dataStack, opcode.Operands)

		fmt.Print("\n")
	}

	return instructions
}
