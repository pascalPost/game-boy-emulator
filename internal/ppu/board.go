package ppu

import (
	"errors"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	start              = float32(-1.0)
	length             = float32(2.0)
	dimensions         = 2
	nTrianglesPerCell  = 2
	nPointsPerTriangle = 3
)

type pixelDisplay struct {
	rows, cols        uint
	points            []float32
	colors            []uint32
	vertexArrayObject uint32
	vboVertices       uint32
	vboColors         uint32
}

func (d *pixelDisplay) createBuffers() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Vertex buffer
	var vboVertices uint32
	gl.GenBuffers(1, &vboVertices)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboVertices)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(d.points), gl.Ptr(d.points), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	// Color buffer
	var vboColors uint32
	gl.GenBuffers(1, &vboColors)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboColors)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(d.colors), gl.Ptr(d.colors), gl.DYNAMIC_DRAW)
	gl.VertexAttribIPointer(1, 1, gl.UNSIGNED_INT, 0, nil)
	gl.EnableVertexAttribArray(1)

	d.vertexArrayObject = vao
	d.vboVertices = vboVertices
	d.vboColors = vboColors
}

func newDisplay(showPixel bool, rows, cols uint) *pixelDisplay {
	d := &pixelDisplay{rows: rows, cols: cols}

	pixelDistance := float32(0.0)

	if showPixel {
		pixelDistance = 0.002
	}

	nCells := rows * cols
	nPoints := nCells * nTrianglesPerCell * nPointsPerTriangle

	dx := length / float32(rows)
	dy := length / float32(cols)

	d.points = make([]float32, nPoints*dimensions)

	for j := uint(0); j < cols; j++ {
		yLow := start + dy*float32(j) + pixelDistance
		yUp := start + dy*float32(j) + dy - pixelDistance

		for i := uint(0); i < rows; i++ {
			xLeft := start + dx*float32(i) + pixelDistance
			xRight := start + dx*float32(i) + dx - pixelDistance

			square := []float32{
				xLeft, yUp,
				xLeft, yLow,
				xRight, yLow,

				xLeft, yUp,
				xRight, yUp,
				xRight, yLow,
			}

			index := uint(len(square)) * (i + j*rows)

			copy(d.points[index:index+uint(len(square))], square)
		}
	}

	d.colors = make([]uint32, nPoints)

	d.createBuffers()

	return d
}

func (d *pixelDisplay) draw() {
	gl.BindVertexArray(d.vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(d.points)/dimensions))
}

func (d *pixelDisplay) updateColors(colors []uint8) error {
	// cols, rows
	if d.rows*d.cols != uint(len(colors)) {
		return errors.New("number of colors does not match number of cells")
	}

	for i, color := range colors {
		for cellPoint := 0; cellPoint < nTrianglesPerCell*nPointsPerTriangle; cellPoint++ {
			d.colors[i*nTrianglesPerCell*nPointsPerTriangle+cellPoint] = uint32(color)
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, d.vboColors)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(d.colors), gl.Ptr(d.colors))

	return nil
}
