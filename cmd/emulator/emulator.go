package main

import (
	"flag"
	"fmt"
	"github.com/pascalPost/game-boy-emulator/internal"
	"log"
	"os"
)

func parseArguments() (string, uint16) {
	help := flag.Bool("h", false, "Display this help message and exit")
	startAddress := flag.Uint("start", 0x0100, "Specify the address to start from (defaults to 0x0100)")

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
	fmt.Println("Usage: emulator [OPTIONS] FILE")
	fmt.Println("Options:")
	fmt.Println("  -h       Display this help message and exit")
	fmt.Println("  -start   Specify the start address")
}

func main() {
	fileName, startAddress := parseArguments()
	gb := internal.NewGameBoy()
	err := gb.LoadCartridge(fileName)
	if err != nil {
		log.Panicf("error loading cartridge: %v", err)
	}
	gb.Run(startAddress)
}
