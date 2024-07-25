package internal

import (
	"log/slog"
	"os"
)

type GameBoy struct {
	cpu    cpu
	memory memory
}

func (gb *GameBoy) LoadCartridge(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error in reading rom", "error", err)
		return err
	}

	copy(gb.memory.data[0:], rom)

	return nil
}

func NewGameBoy() *GameBoy {
	return &GameBoy{}
}

func (gb *GameBoy) Run(startAddress uint16) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	const initialStackPointerAddress uint16 = 0xFFFE
	gb.cpu.registers.pc = startAddress
	gb.cpu.registers.sp = initialStackPointerAddress

	for {
		gb.cpu.runInstruction(&gb.memory)
	}
}
