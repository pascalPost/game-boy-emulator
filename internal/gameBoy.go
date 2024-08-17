package internal

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions"
	"github.com/pascalPost/game-boy-emulator/internal/ppu"
	"log/slog"
	"os"
	"runtime"
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

	// copy logo to vram

	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}

	for i, b := range logo {
		gb.Memory.Write(uint16(0x104+i), b)
	}

	// init graphics

	runtime.LockOSThread()

	ppu.InitGlfw()
	defer glfw.Terminate()

	ppu.InitOpenGL()

	const pixelSize = 4
	screenWindow := ppu.CreateWindow(160*pixelSize, 144*pixelSize, "game-boy-emulator", nil)
	ppu.ActivateDebugOutput()

	bgMapWindow := ppu.CreateWindow(750, 750, "BGMap", screenWindow)
	ppu.ActivateDebugOutput()

	bgMapWindow.MakeContextCurrent()

	pixelData := ppu.InitBGMapPixels()

	bgMapVertexColorData := make([]uint32, pixelData.NumVertices)

	// Color buffer
	var vboColors uint32
	gl.GenBuffers(1, &vboColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(bgMapVertexColorData), gl.Ptr(bgMapVertexColorData), gl.DYNAMIC_DRAW)
	gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	bgMapPixelColorData := make([]uint8, 0, 32*32*8*8)

	viewportVertexData := ppu.UpdateViewport(&gb.Memory)
	viewportData := ppu.InitViewport(viewportVertexData)

	gridData := ppu.InitGrid(ppu.BgMapRows, ppu.BgMapCols)

	const initialStackPointerAddress uint16 = 0xFFFE
	gb.Cpu.Registers.PC = startAddress
	gb.Cpu.Registers.SP = initialStackPointerAddress

	const framerate = 30
	lastTime := glfw.GetTime()

	// hack LY init; needs to be replaced with proper LY updates
	gb.Memory.Write(0xFF44, 0x90)

	for !screenWindow.ShouldClose() && !bgMapWindow.ShouldClose() {

		instructions.RunInstruction(&gb.Cpu, &gb.Memory)
		//}

		if gb.Cpu.Registers.PC == uint16(0x0066) {
			println("Test")
		}

		currentTime := glfw.GetTime()
		if currentTime-lastTime > 1.0/framerate {
			lastTime = currentTime

			screenWindow.MakeContextCurrent()
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			screenWindow.SwapBuffers()

			bgMapWindow.MakeContextCurrent()

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			bgMapPixelColorData = bgMapPixelColorData[0:0:cap(bgMapPixelColorData)]
			// assuming 0x9800 tile map and 0x8000 tile addressing
			const tileBasePointer = 0x8000
			for vRamAddress := uint16(0x9800); vRamAddress < 0x9C00; vRamAddress++ {
				tileIndex := gb.Memory.Read(vRamAddress)
				tileAddress := tileBasePointer + int(tileIndex)*ppu.TileByteSize
				tileData := gb.Memory.Data[tileAddress : tileAddress+ppu.TileByteSize]
				bgMapPixelColorData = ppu.ConvertIntoPixelColors(tileData, bgMapPixelColorData)
			}

			for i, color := range bgMapPixelColorData {
				for cellPoint := 0; cellPoint < ppu.NumTrianglesPerCell*ppu.NumPointsPerTriangle; cellPoint++ {
					bgMapVertexColorData[i*ppu.NumTrianglesPerCell*ppu.NumPointsPerTriangle+cellPoint] = uint32(color)
				}
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, ppu.Float32ByteSize*len(bgMapVertexColorData), gl.Ptr(bgMapVertexColorData))

			gl.UseProgram(pixelData.Program)
			gl.BindVertexArray(pixelData.VertexArrayObject)
			gl.DrawArrays(gl.TRIANGLES, 0, pixelData.NumVertices)

			gl.UseProgram(gridData.Program)
			gl.BindVertexArray(gridData.VertexArrayObject)
			gl.DrawArrays(gl.LINES, 0, gridData.NumVertices)

			viewportVertexData := ppu.UpdateViewport(&gb.Memory)
			gl.BindBuffer(gl.ARRAY_BUFFER, viewportData.VertexBufferObject)
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, ppu.Float32ByteSize*len(viewportVertexData), gl.Ptr(viewportVertexData[:]))

			gl.UseProgram(viewportData.Program)
			gl.BindVertexArray(viewportData.VertexArrayObject)
			gl.DrawArrays(gl.LINES, 0, viewportData.NumVertices)

			bgMapWindow.SwapBuffers()

			glfw.PollEvents()
		}
	}
}
