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

func (gb *GameBoy) Run() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	const headerEntryAddress uint16 = 0x0100
	const initialStackPointerAddress uint16 = 0xFFFE
	gb.cpu.registers.pc = headerEntryAddress
	gb.cpu.registers.sp = initialStackPointerAddress

	for {
		gb.cpu.runInstruction(&gb.memory)
	}
}
