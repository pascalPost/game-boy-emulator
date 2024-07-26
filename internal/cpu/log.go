package cpu

import (
	"fmt"
	"log/slog"
)

func FmtHex16(value uint16) string {
	return fmt.Sprintf("0x%04X", value)
}

func FmtHex8(value uint8) string {
	return fmt.Sprintf("0x%02X", value)
}

func Log(memory *Memory, programCounter uint16, instructionLengthInBytes int, instruction, description string) {
	pcBegin := programCounter - 1

	// handle prefixed instruction
	if pcBegin > 0 && memory.Data[pcBegin-1] == uint8(0xCB) {
		pcBegin--
	}

	pcEnd := pcBegin + uint16(instructionLengthInBytes)
	slog.Debug("Instruction", "PC", FmtHex16(pcBegin), "mem", fmt.Sprintf("0x% 2X", memory.Data[pcBegin:pcEnd]), "instruction", instruction, "description", description)
}
