package ppu

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pascalPost/game-boy-emulator/internal/cpu"
	"runtime"
)

const (
	tileRows      = 8
	tileCols      = 8
	tileAreas     = 3
	tilesDataRows = 8
	tilesDataCols = 16
	BgMapRows     = 32
	BgMapCols     = 32
	TileByteSize  = 16
)

func ConvertIntoPixelColors(tile []byte, pixelColor []byte) []byte {
	const tileSize = 16

	if len(tile) < tileSize {
		panic("tile must be 16 bytes long")
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

			pixelColor = append(pixelColor, colorValue)
		}
	}

	return pixelColor
}

func ComputePixelColors(tile []byte) []byte {
	if len(tile) < 16 {
		panic("tile must be 16 bytes long")
	}

	colorValues := make([]byte, 0, tileRows*tileCols)

	colorValues = ConvertIntoPixelColors(tile, colorValues)

	return colorValues
}

func PlotTile(cellColors []uint8) {
	runtime.LockOSThread()

	InitGlfw()
	defer glfw.Terminate()

	window := CreateWindow(500, 500, "tile", nil)

	InitOpenGL()

	program := NewProgram(VertexShader2DColor, FragmentShaderColor)

	var vao uint32
	var vertexDataLength int32
	var vboColors uint32
	{
		const nCells = tileCols * tileCols
		const nPoints = nCells * NumTrianglesPerCell * NumPointsPerTriangle

		points := make([]float32, 0, nPoints*Dimensions)

		points = appendTilePoints(start, start+length, start+length, start, points)

		vertexDataLength = int32(len(points))

		colors := make([]uint32, nPoints)
		for i, color := range cellColors {
			for cellPoint := 0; cellPoint < NumTrianglesPerCell*NumPointsPerTriangle; cellPoint++ {
				colors[i*NumTrianglesPerCell*NumPointsPerTriangle+cellPoint] = uint32(color)
			}
		}

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		// Vertex buffer
		var vboVertices uint32
		gl.GenBuffers(1, &vboVertices)
		gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
		gl.EnableVertexAttribArray(0)

		// Color buffer
		gl.GenBuffers(1, &vboColors)
		gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.DYNAMIC_DRAW)
		gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
		gl.EnableVertexAttribArray(1)
	}

	//_ = display.updateColors(cellColors)

	//for i, color := range colors {
	//	for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
	//		d.colors[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
	//	}
	//}

	//gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(d.colors), gl.Ptr(d.colors))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, vertexDataLength/Dimensions)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PlotBGMap(pixelColorData []byte, top, left, bottom, right uint8) {
	runtime.LockOSThread()

	InitGlfw()
	defer glfw.Terminate()

	window := CreateWindow(750, 750, "BGMap", nil)

	InitOpenGL()

	pixelData := InitBGMapPixels()

	colorData := make([]uint32, pixelData.NumVertices)

	// Color buffer
	var vboColors uint32
	gl.GenBuffers(1, &vboColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colorData), gl.Ptr(colorData), gl.DYNAMIC_DRAW)
	gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	//UpdateViewport()
	//viewportData := InitViewport(top, left, bottom, right)

	gridData := InitGrid(BgMapRows, BgMapCols)

	for i, color := range pixelColorData {
		for cellPoint := 0; cellPoint < NumTrianglesPerCell*NumPointsPerTriangle; cellPoint++ {
			colorData[i*NumTrianglesPerCell*NumPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(colorData), gl.Ptr(colorData))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(pixelData.Program)
		gl.BindVertexArray(pixelData.VertexArrayObject)
		gl.DrawArrays(gl.TRIANGLES, 0, pixelData.NumVertices)

		gl.UseProgram(gridData.Program)
		gl.BindVertexArray(gridData.VertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, gridData.NumVertices)

		//gl.UseProgram(viewportData.Program)
		//gl.BindVertexArray(viewportData.VertexArrayObject)
		//gl.DrawArrays(gl.LINES, 0, viewportData.NumVertices)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PlotTiles(pixelData []byte) {
	runtime.LockOSThread()

	InitGlfw()
	defer glfw.Terminate()

	window := CreateWindow(500, 750, "tiles", nil)

	InitOpenGL()

	tilePixelData := initTilesPixels()

	colorData := make([]uint32, tilePixelData.NumVertices*Dimensions)

	// Color buffer
	var vboColors uint32
	gl.GenBuffers(1, &vboColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colorData), gl.Ptr(colorData), gl.DYNAMIC_DRAW)
	gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	tileMapGridData := initTilesGrid()

	tileMapSplitData := initTilesSplit()

	for i, color := range pixelData {
		for cellPoint := 0; cellPoint < NumTrianglesPerCell*NumPointsPerTriangle; cellPoint++ {
			colorData[i*NumTrianglesPerCell*NumPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(colorData), gl.Ptr(colorData))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(tilePixelData.Program)
		gl.BindVertexArray(tilePixelData.VertexArrayObject)
		gl.DrawArrays(gl.TRIANGLES, 0, tilePixelData.NumVertices)

		gl.UseProgram(tileMapGridData.Program)
		gl.BindVertexArray(tileMapGridData.VertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, tileMapGridData.NumVertices)

		gl.UseProgram(tileMapSplitData.Program)
		gl.BindVertexArray(tileMapSplitData.VertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, tileMapSplitData.NumVertices)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

type GlData struct {
	Program            uint32
	VertexArrayObject  uint32
	NumVertices        int32
	VertexBufferObject uint32
}

func initTilesSplit() GlData {
	tileMapSplitProgram := NewProgram(vertexShader2DNoColor, fragmentShaderRed)

	tileMapSplitLineData := []float32{-1.0, -0.333333, 1.0, -0.333333, -1.0, 0.333333, 1.0, 0.333333}
	tileMapSplitVertices := int32(len(tileMapSplitLineData)) / 2

	var tileMapSplitVao uint32
	gl.GenVertexArrays(1, &tileMapSplitVao)
	gl.BindVertexArray(tileMapSplitVao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(tileMapSplitLineData), gl.Ptr(tileMapSplitLineData), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	return GlData{tileMapSplitProgram, tileMapSplitVao, tileMapSplitVertices, vboVertices}
}

func initTilesGrid() GlData {
	tileMapGridProgram := NewProgram(vertexShader2DNoColor, fragmentShaderBlack)

	const deltaX = length / tilesDataCols
	const deltaY = length / tileAreas / tilesDataRows

	const nVertices = (tilesDataCols-1)*2 + (tilesDataRows-1)*2*tileAreas

	tileMapGridLineData := make([]float32, 0, nVertices)

	const end = length - start
	const tol = 0.001

	const endX = end - deltaX + tol
	for x := start + deltaX; x < endX; x += deltaX {
		tileMapGridLineData = append(tileMapGridLineData, x, start, x, end)
	}

	const endY = end - deltaY + tol
	for y := start + deltaY; y < endY; y += deltaY {
		tileMapGridLineData = append(tileMapGridLineData, start, y, end, y)
	}

	var tileMapGridVertices int32
	tileMapGridVertices = int32(len(tileMapGridLineData)) / 2

	var tileMapGridVao uint32
	gl.GenVertexArrays(1, &tileMapGridVao)
	gl.BindVertexArray(tileMapGridVao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(tileMapGridLineData), gl.Ptr(tileMapGridLineData), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return GlData{tileMapGridProgram, tileMapGridVao, tileMapGridVertices, vboVertices}
}

func InitGrid(rows, cols uint) GlData {
	// TODO merge with initTilesGrid

	tileMapGridProgram := NewProgram(vertexShader2DNoColor, fragmentShaderBlack)

	deltaX := length / float32(cols)
	deltaY := length / float32(rows)

	nVertices := (cols-1)*2 + (rows-1)*2

	tileMapGridLineData := make([]float32, 0, nVertices)

	const end = length - start
	const tol = 0.001

	endX := end - deltaX + tol
	for x := start + deltaX; x < endX; x += deltaX {
		tileMapGridLineData = append(tileMapGridLineData, x, start, x, end)
	}

	endY := end - deltaY + tol
	for y := start + deltaY; y < endY; y += deltaY {
		tileMapGridLineData = append(tileMapGridLineData, start, y, end, y)
	}

	var tileMapGridVertices int32
	tileMapGridVertices = int32(len(tileMapGridLineData)) / 2

	var tileMapGridVao uint32
	gl.GenVertexArrays(1, &tileMapGridVao)
	gl.BindVertexArray(tileMapGridVao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(tileMapGridLineData), gl.Ptr(tileMapGridLineData), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return GlData{tileMapGridProgram, tileMapGridVao, tileMapGridVertices, vboVertices}
}

func appendTilePoints(xTileStart, xTileEnd, yTileStart, yTileEnd float32, data []float32) []float32 {
	xTileLength := xTileEnd - xTileStart
	yTileLength := yTileEnd - yTileStart

	deltaX := xTileLength / float32(tileCols)
	deltaY := yTileLength / float32(tileRows)

	// xStart    xEnd
	// ------- x yStart
	// |
	// |
	// |
	// y		 yEnd

	for rowIdx := 0; rowIdx < tileRows; rowIdx++ {
		yCellStart := yTileStart + deltaY*float32(rowIdx)
		yCellEnd := yCellStart + deltaY

		for colIdx := 0; colIdx < tileCols; colIdx++ {
			xCellStart := xTileStart + deltaX*float32(colIdx)
			xCellEnd := xCellStart + deltaX

			// square
			data = append(data,
				xCellStart, yCellEnd,
				xCellStart, yCellStart,
				xCellEnd, yCellStart,
				xCellStart, yCellEnd,
				xCellEnd, yCellEnd,
				xCellEnd, yCellStart)
		}
	}

	return data
}

func InitBGMapPixels() GlData {
	tileMapGridProgram := NewProgram(VertexShader2DColor, FragmentShaderColor)

	const deltaX = length / BgMapCols
	const deltaY = length / BgMapRows

	const nTiles = BgMapCols * BgMapRows
	const nPixels = nTiles * tileRows * tileCols

	points := make([]float32, 0, nPixels*NumTrianglesPerCell*NumPointsPerTriangle*Dimensions)

	for rowIdx := 0; rowIdx < BgMapRows; rowIdx++ {
		yStart := start + length - deltaY*float32(rowIdx)
		yEnd := yStart - deltaY

		for colIdx := 0; colIdx < BgMapCols; colIdx++ {
			xStart := start + deltaX*float32(colIdx)
			xEnd := xStart + deltaX

			points = appendTilePoints(xStart, xEnd, yStart, yEnd, points)
		}
	}

	var nVertices int32
	nVertices = int32(len(points)) / Dimensions

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return GlData{tileMapGridProgram, vao, nVertices, vboVertices}
}

func initTilesPixels() GlData {
	tileMapGridProgram := NewProgram(VertexShader2DColor, FragmentShaderColor)

	const rows = tileAreas * tilesDataRows
	const cols = tilesDataCols

	const deltaX = length / cols
	const deltaY = length / rows

	const nCells = rows * cols

	points := make([]float32, 0, nCells*NumTrianglesPerCell*NumPointsPerTriangle*Dimensions)

	for mapIdx := 0; mapIdx < tileAreas; mapIdx++ {
		yMapStart := start + length - length/tileAreas*float32(mapIdx)
		//yMapEnd := yMapStart - length/tileAreas

		for rowIdx := 0; rowIdx < tilesDataRows; rowIdx++ {
			yStart := yMapStart - deltaY*float32(rowIdx)
			yEnd := yStart - deltaY

			for colIdx := 0; colIdx < tilesDataCols; colIdx++ {
				xStart := start + deltaX*float32(colIdx)
				xEnd := xStart + deltaX

				// create tile pixels
				points = appendTilePoints(xStart, xEnd, yStart, yEnd, points)
			}
		}
	}

	var nVertices int32
	nVertices = int32(len(points)) / Dimensions

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return GlData{tileMapGridProgram, vao, nVertices, vboVertices}
}

func viewportPositionTransformX(value uint8) float32 {
	return start + float32(value)*length/255
}

func viewportPositionTransformY(value uint8) float32 {
	return end - float32(value)*length/255
}

func InitViewport(vertexData *[16]float32) GlData {

	program := NewProgram(vertexShader2DNoColor, fragmentShaderRed)

	nVertices := int32(len(vertexData)) / Dimensions

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// TODO introduce vbo to reduce mem footprint

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertexData), gl.Ptr(vertexData[:]), gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(0, Dimensions, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	return GlData{program, vao, nVertices, vboVertices}
}

func UpdateViewport(memory *cpu.Memory) *[16]float32 {
	const scyAddress = 0xFF42
	const scxAddress = 0xFF43

	top := memory.Read(scyAddress)
	left := memory.Read(scxAddress)

	bottom := uint8((uint16(top) + 143) % 256)
	right := uint8((uint16(left) + 159) % 256)

	return &[16]float32{viewportPositionTransformX(left), viewportPositionTransformY(top),
		viewportPositionTransformX(right), viewportPositionTransformY(top),
		viewportPositionTransformX(left), viewportPositionTransformY(top),
		viewportPositionTransformX(left), viewportPositionTransformY(bottom),
		viewportPositionTransformX(left), viewportPositionTransformY(bottom),
		viewportPositionTransformX(right), viewportPositionTransformY(bottom),
		viewportPositionTransformX(right), viewportPositionTransformY(top),
		viewportPositionTransformX(right), viewportPositionTransformY(bottom)}
}
