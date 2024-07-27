package internal

import (
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions"
	"log/slog"
	"os"
)

type GameBoy struct {
	Cpu    cpu.Cpu
	Memory cpu.Memory
}

func (gb *GameBoy) LoadCartridge(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error in reading rom", "error", err)
		return err
	}

	copy(gb.Memory.Data[0:], rom)

	return nil
}

func NewGameBoy() *GameBoy {
	return &GameBoy{}
}

func (gb *GameBoy) Run(startAddress uint16) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	const initialStackPointerAddress uint16 = 0xFFFE
	gb.Cpu.Registers.PC = startAddress
	gb.Cpu.Registers.SP = initialStackPointerAddress

	for {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}
}
