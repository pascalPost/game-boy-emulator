package main

import (
	"flag"
	"fmt"
	"github.com/pascalPost/game-boy-emulator/internal"
	"io"
	"log"
	"log/slog"
	"os"
)

func parseArguments() (string, uint16) {
	help := flag.Bool("h", false, "Display this help message and exit")
	startAddress := flag.Uint("start", 0x0150, "Specify the address to start from (defaults to 0x0150)")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		printHelp()
		os.Exit(1)
	}

	return flag.Arg(0), uint16(*startAddress)
}

func printHelp() {
	fmt.Println("Usage: disassembler [OPTIONS] FILE")
	fmt.Println("Options:")
	fmt.Println("  -h       Display this help message and exit")
	fmt.Println("  -start   Specify the start address")
}

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
	fileName, startAddress := parseArguments()

	opcodes, err := internal.ParseOpcodes()
	if err != nil {
		log.Fatal(err)
	}

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

	instructions := internal.Disassemble(data, startAddress, opcodes)
	printInstructions(data, instructions, 0)
}
