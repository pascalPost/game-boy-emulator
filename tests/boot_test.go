package tests

import (
	"github.com/pascalPost/game-boy-emulator/internal"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions"
	"github.com/pascalPost/game-boy-emulator/internal/ppu"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoot(t *testing.T) {
	gb := internal.NewGameBoy()
	err := gb.LoadCartridge("../roms/DMG_ROM.bin")
	assert.NoError(t, err)

	gb.Cpu.Registers.PC = 0x0000

	// copy nintendo logo to vram

	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}

	for i, b := range logo {
		gb.Memory.Write(uint16(0x104+i), b)
	}

	// setup stack

	// 0x31 FE FF : LD SP, 0xFFFE
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0003), gb.Cpu.Registers.PC)
	assert.Equal(t, uint16(0xFFFE), gb.Cpu.Registers.SP)

	// zero memory from 0x8000 to 0x9FFF (VRAM)

	// 0xAF : XOR A, A
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0004), gb.Cpu.Registers.PC)
	assert.Equal(t, uint8(0x00), gb.Cpu.Registers.A())
	assert.True(t, gb.Cpu.Registers.Flags().Z())
	assert.False(t, gb.Cpu.Registers.Flags().N())
	assert.False(t, gb.Cpu.Registers.Flags().H())
	assert.False(t, gb.Cpu.Registers.Flags().C())

	// 0x21 00 00 : LD HL, 0x9FFFF
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0007), gb.Cpu.Registers.PC)
	assert.Equal(t, uint16(0x9FFF), gb.Cpu.Registers.HL)

	for i := uint16(0x9FFF); i > 0x8000; i-- {
		// 0x32 : LDD [HL], A
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
		assert.Equal(t, uint16(0x0008), gb.Cpu.Registers.PC)
		assert.Equal(t, i-1, gb.Cpu.Registers.HL)
		assert.Equal(t, uint8(0x00), gb.Memory.Read(i))

		// 0xCB 7C : BIT 7, H
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
		assert.Equal(t, uint16(0x000A), gb.Cpu.Registers.PC)
		assert.False(t, gb.Cpu.Registers.Flags().Z())
		assert.False(t, gb.Cpu.Registers.Flags().N())
		assert.True(t, gb.Cpu.Registers.Flags().H())
		assert.False(t, gb.Cpu.Registers.Flags().C(), "C flag is unaffected and thus still false")

		// 0x20 FB : JR NZ, -5
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
		assert.Equal(t, uint16(0x0007), gb.Cpu.Registers.PC)
	}

	// 0x32 : LDD [HL], A
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0008), gb.Cpu.Registers.PC)
	assert.Equal(t, uint16(0x7FFF), gb.Cpu.Registers.HL)
	assert.Equal(t, uint8(0x00), gb.Memory.Read(0x8000))

	// 0xCB 7C : BIT 7, H
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x000A), gb.Cpu.Registers.PC)
	assert.True(t, gb.Cpu.Registers.Flags().Z())
	assert.False(t, gb.Cpu.Registers.Flags().N())
	assert.True(t, gb.Cpu.Registers.Flags().H())
	assert.False(t, gb.Cpu.Registers.Flags().C(), "C flag is unaffected and thus still false")

	// 0x20 FB : JR NZ, -5
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x000C), gb.Cpu.Registers.PC)

	// Setup Audio

	// 0x21 26 FF : LD HL, 0xFF26
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x000F), gb.Cpu.Registers.PC)
	assert.Equal(t, uint16(0xFF26), gb.Cpu.Registers.HL)

	// 0x0E 11 : LD C, 0x11
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0011), gb.Cpu.Registers.PC)
	assert.Equal(t, uint8(0x11), gb.Cpu.Registers.C())

	// 0x3E 80 : LD A, 0x80
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0013), gb.Cpu.Registers.PC)
	assert.Equal(t, uint8(0x80), gb.Cpu.Registers.A())

	// 0x32 : LDD [HL], A
	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	assert.Equal(t, uint16(0x0014), gb.Cpu.Registers.PC)
	assert.Equal(t, uint8(0x80), gb.Memory.Read(0xFF26))
	assert.Equal(t, uint16(0xFF25), gb.Cpu.Registers.HL)

	// TODO deviations from bgb

	for gb.Cpu.Registers.PC != uint16(0x0021) {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}

	assert.Equal(t, uint16(0x0021), gb.Cpu.Registers.PC)

	// Convert and load logo data from cartridge into Video RAM (VRAM)

	for gb.Cpu.Registers.PC != uint16(0x0034) {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}

	assert.Equal(t, uint16(0x0034), gb.Cpu.Registers.PC)

	// Load additional bytes into VRAM (the tile for R)

	for gb.Cpu.Registers.PC != uint16(0x0040) {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}

	assert.Equal(t, uint16(0x0040), gb.Cpu.Registers.PC)

	// Setup background tilemap

	for gb.Cpu.Registers.PC != uint16(0x0055) {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}

	assert.Equal(t, uint16(0x0055), gb.Cpu.Registers.PC)

	// Scroll logo on screen, and play logo sound

	for gb.Cpu.Registers.PC != uint16(0x0062) {
		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	}

	assert.Equal(t, uint16(0x0062), gb.Cpu.Registers.PC)

	//fmt.Printf("VRAM: % 02X\n", gb.Memory.Data[0x8000:0x8020])

	//pixels := ppu.ComputePixelColors(gb.Memory.Data[0x8010:0x8020])

	// 2*8*16 tiles with 8*8 each
	// using an uint8 results in a waste of 6 bits; this can be optimized
	tileMapSize := 3 * 8 * 16
	tilePixelSize := 8 * 8
	tileMapPixelData := make([]uint8, 0, tileMapSize*tilePixelSize)

	const tileSize = 16
	const vRAMStart = 0x8000
	for tileIndex := 0; tileIndex < tileMapSize; tileIndex++ {
		tileVRAMStart := vRAMStart + tileIndex*tileSize
		tileVRAMEnd := tileVRAMStart + tileSize
		tileMapPixelData = ppu.ConvertIntoPixelColors(gb.Memory.Data[tileVRAMStart:tileVRAMEnd], tileMapPixelData)
	}

	//ppu.PlotTile(tileMapPixelData[64 : 64+64])
	//ppu.PlotTile(tileMapPixelData[64+64 : 64+64+64])

	ppu.PlotTileMap(tileMapPixelData)

	// TODO plot BGMap
	// TODO plot screen

	//// wait for screen
	//
	//for gb.Cpu.Registers.PC != uint16(0x0068) {
	//	instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	//}
	//
	//assert.Equal(t, uint16(0x0068), gb.Cpu.Registers.PC)
	//
	//slog.SetLogLoggerLevel(slog.LevelDebug)
	//
	//instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	//instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	//instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	//instructions.RunInstruction(&gb.Cpu, &gb.Memory)
	//instructions.RunInstruction(&gb.Cpu, &gb.Memory)

}

//func plotTileMap(cellColors []uint8) {
//	runtime.LockOSThread()
//
//	window := initGlfw()
//	defer glfw.Terminate()
//
//	program := initOpenGL()
//
//	display := newDisplay(true, 8, 8)
//
//	_ = display.updateColors(cellColors)
//
//	for !window.ShouldClose() {
//		draw(display, window, program)
//	}
//}
