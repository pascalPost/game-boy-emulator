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
	program           uint32
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

	const vertexShaderSource = `
		#version 410
		layout(location = 0) in vec2 aPos;
		layout(location = 1) in uint aColor;

		flat out uint color;

		void main() {
			gl_Position = vec4(aPos, 0.0, 1.0);
			color = aColor;
		}
	` + "\x00"

	const fragmentShaderSource = `
		#version 410
		out vec4 FragColor;
		flat in uint color;

		void main() {
			vec3 colorVec = vec3(1.0, 1.0, 1.0);
			if (color == 0u) {
				colorVec = vec3(1.0, 1.0, 1.0); // White
			} else if (color == 1u) {
				colorVec = vec3(0.75, 0.75, 0.75); // Light Gray
			} else if (color == 2u) {
				colorVec = vec3(0.25, 0.25, 0.25); // Dark Gray
			} else if (color == 3u) {
				colorVec = vec3(0.0, 0.0, 0.0); // Black
			} else {
				colorVec = vec3(1.0, 1.0, 0.0); // Default to Yellow
			}

			FragColor = vec4(colorVec, 1.0);
		}
	` + "\x00"

	d.program = NewProgram(vertexShaderSource, fragmentShaderSource)

	pixelDistance := float32(0.0)

	if showPixel {
		pixelDistance = 0.002
	}

	nCells := rows * cols
	nPoints := nCells * nTrianglesPerCell * nPointsPerTriangle

	dx := length / float32(cols)
	dy := length / float32(rows)

	d.points = make([]float32, nPoints*dimensions)

	for j := uint(0); j < rows; j++ {
		yLow := start + dy*float32(j) + pixelDistance
		yUp := start + dy*float32(j) + dy - pixelDistance

		for i := uint(0); i < cols; i++ {
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

			index := uint(len(square)) * (i + j*cols)

			copy(d.points[index:index+uint(len(square))], square)
		}
	}

	d.colors = make([]uint32, nPoints)

	d.createBuffers()

	return d
}

func (d *pixelDisplay) draw() {
	gl.UseProgram(d.program)
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
