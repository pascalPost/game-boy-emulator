package ppu

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

const (
	tileRows      = 8
	tileCols      = 8
	tileAreas     = 3
	tilesDataRows = 8
	tilesDataCols = 16
	bgMapRows     = 32
	bgMapCols     = 32
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

	window := initGlfw(500, 500)
	defer glfw.Terminate()

	initOpenGL()

	program := NewProgram(vertexShader2DColor, fragmentShaderColor)

	var vao uint32
	var vertexDataLength int32
	var vboColors uint32
	{
		const nCells = tileCols * tileCols
		const nPoints = nCells * nTrianglesPerCell * nPointsPerTriangle

		points := make([]float32, 0, nPoints*dimensions)

		points = appendTilePoints(start, start+length, start+length, start, points)

		vertexDataLength = int32(len(points))

		colors := make([]uint32, nPoints)
		for i, color := range cellColors {
			for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
				colors[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
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
		gl.DrawArrays(gl.TRIANGLES, 0, vertexDataLength/dimensions)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PlotBGMap(pixelColorData []byte) {
	runtime.LockOSThread()

	window := initGlfw(750, 750)
	defer glfw.Terminate()

	initOpenGL()

	pixelData := initBGMapPixels()

	colorData := make([]uint32, pixelData.nVertices)

	// Color buffer
	var vboColors uint32
	gl.GenBuffers(1, &vboColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colorData), gl.Ptr(colorData), gl.DYNAMIC_DRAW)
	gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	gridData := initGrid(bgMapRows, bgMapCols)

	for i, color := range pixelColorData {
		for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
			colorData[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(colorData), gl.Ptr(colorData))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(pixelData.program)
		gl.BindVertexArray(pixelData.vertexArrayObject)
		gl.DrawArrays(gl.TRIANGLES, 0, pixelData.nVertices)

		gl.UseProgram(gridData.program)
		gl.BindVertexArray(gridData.vertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, gridData.nVertices)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func PlotTiles(pixelData []byte) {
	runtime.LockOSThread()

	window := initGlfw(500, 750)
	defer glfw.Terminate()

	initOpenGL()

	tilePixelData := initTilesPixels()

	colorData := make([]uint32, tilePixelData.nVertices*dimensions)

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
		for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
			colorData[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(colorData), gl.Ptr(colorData))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(tilePixelData.program)
		gl.BindVertexArray(tilePixelData.vertexArrayObject)
		gl.DrawArrays(gl.TRIANGLES, 0, tilePixelData.nVertices)

		gl.UseProgram(tileMapGridData.program)
		gl.BindVertexArray(tileMapGridData.vertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, tileMapGridData.nVertices)

		gl.UseProgram(tileMapSplitData.program)
		gl.BindVertexArray(tileMapSplitData.vertexArrayObject)
		gl.DrawArrays(gl.LINES, 0, tileMapSplitData.nVertices)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

type glData struct {
	program           uint32
	vertexArrayObject uint32
	nVertices         int32
}

func initTilesSplit() glData {
	tileMapSplitProgram := NewProgram(vertexShader2DNoColor, fragmentShaderRed)

	var tileMapSplitVao uint32
	var tileMapSplitVertices int32
	{
		tileMapSplitLineData := []float32{-1.0, -0.333333, 1.0, -0.333333, -1.0, 0.333333, 1.0, 0.333333}
		tileMapSplitVertices = int32(len(tileMapSplitLineData)) / 2

		gl.GenVertexArrays(1, &tileMapSplitVao)
		gl.BindVertexArray(tileMapSplitVao)

		// Vertex buffer
		var vboVertices uint32
		gl.GenBuffers(1, &vboVertices)
		gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(tileMapSplitLineData), gl.Ptr(tileMapSplitLineData), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
		gl.EnableVertexAttribArray(0)
	}
	return glData{tileMapSplitProgram, tileMapSplitVao, tileMapSplitVertices}
}

func initTilesGrid() glData {
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

	return glData{tileMapGridProgram, tileMapGridVao, tileMapGridVertices}
}

func initGrid(rows, cols uint) glData {
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

	return glData{tileMapGridProgram, tileMapGridVao, tileMapGridVertices}
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

func initBGMapPixels() glData {
	tileMapGridProgram := NewProgram(vertexShader2DColor, fragmentShaderColor)

	const deltaX = length / bgMapCols
	const deltaY = length / bgMapRows

	const nTiles = bgMapCols * bgMapRows
	const nPixels = nTiles * tileRows * tileCols

	points := make([]float32, 0, nPixels*nTrianglesPerCell*nPointsPerTriangle*dimensions)

	for rowIdx := 0; rowIdx < bgMapRows; rowIdx++ {
		yStart := start + length - deltaY*float32(rowIdx)
		yEnd := yStart - deltaY

		for colIdx := 0; colIdx < bgMapCols; colIdx++ {
			xStart := start + deltaX*float32(colIdx)
			xEnd := xStart + deltaX

			points = appendTilePoints(xStart, xEnd, yStart, yEnd, points)
		}
	}

	var nVertices int32
	nVertices = int32(len(points)) / dimensions

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

	return glData{tileMapGridProgram, vao, nVertices}
}

func initTilesPixels() glData {
	tileMapGridProgram := NewProgram(vertexShader2DColor, fragmentShaderColor)

	const rows = tileAreas * tilesDataRows
	const cols = tilesDataCols

	const deltaX = length / cols
	const deltaY = length / rows

	const nCells = rows * cols

	points := make([]float32, 0, nCells*nTrianglesPerCell*nPointsPerTriangle*dimensions)

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
	nVertices = int32(len(points)) / dimensions

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

	return glData{tileMapGridProgram, vao, nVertices}
}
