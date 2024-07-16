package main

import (
	"bufio"
	"fmt"
	"game-boy-emulator/internal"
	"io"
	"log"
	"log/slog"
	"os"
)

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
	internal.Disassemble(dataBuf[:], opcodes)

	fmt.Printf("\n")
	fmt.Printf("Read program (starting from 0x0150):\n")

	//_, err = rom.Seek(0x150, io.SeekStart)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = rom.Read(buf)
	//if err != nil {
	//	log.Fatal(err)
	//}

	_, err = io.ReadAtLeast(br, buf, 20)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("% X\n", buf)
	internal.Disassemble(buf[:], opcodes)
}
