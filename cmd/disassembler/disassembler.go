package main

import (
	"bufio"
	"fmt"
	"game-boy-emulator/internal"
	"io"
	"log"
	"log/slog"
	"os"
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

func readOperands(data []byte, operands []internal.Operand) []byte {
	for _, o := range operands {
		address := toAddress(data[0:o.Bytes])
		fmt.Printf(" 0x%X", address)
		data = data[o.Bytes:]
	}

	return data
}

func disassemble(dataStack []byte, list *internal.OpcodeList) {
	for len(dataStack) > 0 {
		prefixed := isPrefixed(dataStack[0])

		var opcode internal.Opcode
		if !prefixed {
			opcode = list.UnPrefixed[internal.ByteKey{dataStack[0]}]
		} else {
			//	i += 1
			dataStack = dataStack[1:]
			opcode = list.UnPrefixed[internal.ByteKey{dataStack[0]}]
		}
		dataStack = dataStack[1:]

		fmt.Printf("%s", opcode.Mnemonic)

		dataStack = readOperands(dataStack, opcode.Operands)

		fmt.Print("\n")
	}
}

func main() {
	opcodes, err := internal.ParseOpcodes()
	if err != nil {
		log.Fatal(err)
	}

	fileName := "roms/snake.gb"

	rom, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Error in opening rom")
	}
	defer func() {
		err := rom.Close()
		if err != nil {
			slog.Error("error in closing rom file")
		}
	}()

	br := bufio.NewReader(rom)

	buf := make([]byte, 0x0150)
	_, err = io.ReadAtLeast(br, buf, 0x0150)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Header entry point:\n")

	header, err := internal.NewHeader(buf)
	dataBuf := header.Raw.EntryPoint
	fmt.Printf("% X\n", dataBuf)
	disassemble(dataBuf[:], opcodes)

	fmt.Printf("\n")
	fmt.Printf("Read program (starting from 0x0150):\n")

	_, err = io.ReadAtLeast(br, buf, 20)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% X\n", dataBuf)
	disassemble(dataBuf[:], opcodes)
}
