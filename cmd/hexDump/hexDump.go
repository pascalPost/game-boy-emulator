package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

func replaceNonPrintChars(chars []byte, replaceWith byte) []byte {
	for i, c := range chars {
		// all ASCII values below 32 (space) and starting from 127 (DEL) are replaced
		if c < 32 || c > 126 {
			chars[i] = replaceWith
		}
	}
	return chars
}

func printData(data []byte) {
	stepSize := uint(4)

	printLine := func(line uint, data []byte, sep string) {
		fmt.Printf("%04X : % X%s%4s\n", line*stepSize, data, sep, replaceNonPrintChars(data, '.'))
	}

	length := uint(len(data))
	line := uint(0)

	for ; line < length/stepSize; line++ {
		start := line * stepSize
		printLine(line, data[start:start+stepSize], "   ")
	}

	remainder := length % 4
	if remainder > 0 {
		start := length - remainder
		sep := strings.Repeat(" ", (4-int(remainder))*3+int(remainder)-1)
		printLine(line, data[start:], sep)
	}
}

func printHelp() {
	fmt.Println("Usage: hexDump [OPTIONS] FILE")
	fmt.Println("Options:")
	fmt.Println("  -h    Display this help message and exit")
}

func fileNameFromArguments() string {
	help := flag.Bool("h", false, "Display this help message and exit")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		printHelp()
		os.Exit(1)
	}

	return flag.Arg(0)
}

func main() {
	fileName := fileNameFromArguments()

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

	byteData, err := io.ReadAll(rom)
	if err != nil {
		log.Panicf("Error in reading rom")
	}

	printData(byteData)
}
