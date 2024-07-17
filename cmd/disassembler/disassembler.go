package main

import (
	"fmt"
	"game-boy-emulator/internal"
	"io"
	"log"
	"log/slog"
	"os"
)

func printInstructions(data []byte, instructions []internal.Instruction, offset uint16) {
	for _, instruction := range instructions {
		for i := 0; i < int(instruction.AddressEnd-instruction.AddressStart); i++ {
			address := instruction.AddressStart + uint16(i) + offset
			if i == 0 {
				fmt.Printf("0x%04X %02X %s\n", address, data[address], instruction.Line)
			} else {
				fmt.Printf("0x%04X %02X\n", address, data[address])
			}
		}
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

	data, err := io.ReadAll(rom)
	if err != nil {
		slog.Error("error in reading rom")
	}

	fmt.Printf("Header entry point:\n")

	header, err := internal.NewHeader(data)
	if err != nil {
		log.Panicf("error on reading header: %s", err)
	}
	entryOffset := uint16(0x0100)
	entry := header.Raw.EntryPoint
	instructions := internal.Disassemble(entry[:], 0, opcodes)
	printInstructions(data, instructions, entryOffset)

	fmt.Printf("\n")
	fmt.Printf("Read program:\n")
	instructions = internal.Disassemble(data[:0x170], 0x0150, opcodes)
	printInstructions(data, instructions, 0)
}
