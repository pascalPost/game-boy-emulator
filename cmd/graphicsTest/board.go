package main

import "github.com/go-gl/gl/v4.1-core/gl"

const (
	rows    = 160
	columns = 144
	start   = -1.0
	length  = 2.0
	dx      = length / rows
	dy      = length / columns
)

type display struct {
	points            []float32
	vertexArrayObject uint32
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {

	// TODO enhance with EBO to avoid too much points

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.DYNAMIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func newDisplay() *display {
	d := display{}

	pixelDistance := float32(0.0)

	showPixel := true
	if showPixel {
		pixelDistance = 0.001
	}

	const dimensions = 3

	const nCells = rows * columns
	const nTrianglesPerCell = 2
	const nPointsPerTriangle = 3
	const nPoints = nCells * nTrianglesPerCell * nPointsPerTriangle

	d.points = make([]float32, nPoints*dimensions)

	for j := 0; j < columns; j++ {
		yLow := start + dy*float32(j) - pixelDistance
		yUp := yLow + dy - pixelDistance

		for i := 0; i < rows; i++ {
			xLeft := start + dx*float32(i) - pixelDistance
			xRight := xLeft + dx - pixelDistance

			square := []float32{
				xLeft, yUp, 0,
				xLeft, yLow, 0,
				xRight, yLow, 0,

				xLeft, yUp, 0,
				xRight, yUp, 0,
				xRight, yLow, 0,
			}

			index := len(square) * (i + j*rows)

			copy(d.points[index:index+len(square)], square)
		}
	}

	d.vertexArrayObject = makeVao(d.points)

	return &d
}

func (d *display) draw() {
	gl.BindVertexArray(d.vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(d.points)/3))
}

//func newCell(x, y int) *cell {
//	points := make([]float32, len(square), len(square))
//	copy(points, square)
//
//	for i := 0; i < len(points); i++ {
//		var position float32
//		var size float32
//		switch i % 3 {
//		case 0:
//			size = 1.0 / float32(columns)
//			position = float32(x) * size
//		case 1:
//			size = 1.0 / float32(rows)
//			position = float32(y) * size
//		default:
//			continue
//		}
//
//		if points[i] < 0 {
//			points[i] = (position * 2) - 1
//		} else {
//			points[i] = ((position + size) * 2) - 1
//		}
//	}
//
//	return &cell{
//		drawable: makeVao(points),
//
//		x: x,
//		y: y,
//	}
//}
//
//func makeCells() [][]*cell {
//	cells := make([][]*cell, rows, columns)
//	for x := 0; x < rows; x++ {
//		for y := 0; y < columns; y++ {
//			c := newCell(x, y)
//			cells[x] = append(cells[x], c)
//		}
//	}
//
//	return cells
//}
