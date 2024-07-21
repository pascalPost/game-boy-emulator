package main

import (
	"github.com/pascalPost/game-boy-emulator/cmd"
	"github.com/pascalPost/game-boy-emulator/internal"
	"log"
)

func main() {
	fileName := cmd.FileNameFromArguments("emulator")
	gb := internal.NewGameBoy()
	err := gb.LoadCartridge(fileName)
	if err != nil {
		log.Panicf("error loading cartridge: %v", err)
	}
	gb.Run()
}
