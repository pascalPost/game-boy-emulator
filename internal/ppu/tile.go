package ppu

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"slices"
)

func ConvertIntoPixelColors(tile []byte, pixelColor []byte) {
	const tileSize = 16

	if len(tile) < tileSize {
		panic("tile must be 16 bytes long")
	}

	const cols = 8
	const rows = 8

	if len(pixelColor) < cols*rows {
		panic("pixel color must be at least 64 bytes long (8*8) pixel color values")
	}

	const bytePerRow = 2
	const byteBitSize = 8
	const lastBitIndexInByte = byteBitSize - 1

	for tileByteIndex := byte(0); tileByteIndex < tileSize; tileByteIndex += bytePerRow {
		byte1 := tile[tileByteIndex]
		byte2 := tile[tileByteIndex+1]

		// loop over individual bits
		for b := byte(0); b < byteBitSize; b++ {
			mask := byte(0b0000_0001)

			bit1 := (byte1 >> (lastBitIndexInByte - b)) & mask
			bit2 := (byte2 >> (lastBitIndexInByte - b)) & mask

			colorValue := (bit2 << 1) | bit1

			rowIndex := lastBitIndexInByte - tileByteIndex/2

			colorIndex := rowIndex*cols + b

			pixelColor[colorIndex] = colorValue
		}
	}
}

func ComputePixelColors(tile []byte) []byte {
	if len(tile) < 16 {
		panic("tile must be 16 bytes long")
	}

	const cols = 8
	const rows = 8
	colorValues := make([]byte, cols*rows)

	ConvertIntoPixelColors(tile, colorValues)

	return colorValues
}

func PlotTile(cellColors []uint8) {
	runtime.LockOSThread()

	window := initGlfw(500, 500)
	defer glfw.Terminate()

	initOpenGL()

	display := newDisplay(true, 8, 8)

	_ = display.updateColors(cellColors)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		display.draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PlotTileMap(pixelData []byte) {
	runtime.LockOSThread()

	window := initGlfw(500, 750)
	defer glfw.Terminate()

	initOpenGL()

	display := newDisplay(false, 3*8*8, 16*8)

	//linesProgram := NewProgram(vertexShaderDefaultColorSource, fragmentShaderDefaultColorSource)

	//line := []float32{-1.0, -0.333333, 1.0, -0.333333, -1.0, 0.333333, 1.0, 0.333333}

	//var vao uint32
	//gl.GenVertexArrays(1, &vao)
	//gl.BindVertexArray(vao)
	//
	//// Vertex buffer
	//var vboVertices uint32
	//gl.GenBuffers(1, &vboVertices)
	//gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	//gl.BufferData(gl.ARRAY_BUFFER, 4*len(line), gl.Ptr(line), gl.STATIC_DRAW)
	//gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	//gl.EnableVertexAttribArray(0)

	//_ = display.updateColors(pixelData)

	// first tile to upper left corner
	//pixelDataSlice := pixelData[64 : 64+64]

	for i, color := range pixelData {
		for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
			display.colors[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	slices.Reverse(display.colors)

	gl.BindBuffer(gl.ARRAY_BUFFER, display.vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(display.colors), gl.Ptr(display.colors))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		display.draw()

		//gl.UseProgram(linesProgram)
		//gl.BindVertexArray(vao)
		//gl.DrawArrays(gl.LINES, 0, 4)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
