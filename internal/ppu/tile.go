package ppu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

func GetPixels(tile []byte) [8][8]byte {
	if len(tile) < 16 {
		panic("tile must be 16 bytes long")
	}

	colorValues := [8][8]byte{}

	for i := 0; i < 8; i++ {
		start := i * 2
		byte1 := tile[start]
		byte2 := tile[start+1]

		for b := byte(0); b < 8; b++ {
			mask := byte(0b0000_0001)

			bit1 := (byte1 >> (7 - b)) & mask
			bit2 := (byte2 >> (7 - b)) & mask

			colorValue := (bit2 << 1) | bit1
			colorValues[i][b] = colorValue
		}
	}

	return colorValues
}

func PlotTile(cellColors []uint8) {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	display := newDisplay(true, 8, 8)

	_ = display.updateColors(cellColors)

	for !window.ShouldClose() {
		draw(display, window, program)
	}
}
