package ppu

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"strings"
)

const (
	start              = float32(-1.0)
	length             = float32(2.0)
	dimensions         = 2
	nTrianglesPerCell  = 2
	nPointsPerTriangle = 3

	vertexShader2DNoColor = `
		#version 410
		layout(location = 0) in vec2 aPos;

		void main() {
			gl_Position = vec4(aPos, 0.0, 1.0);
		}
	` + "\x00"

	vertexShader2DColor = `
		#version 410
		layout(location = 0) in vec2 aPos;
		layout(location = 1) in uint aColor;

		flat out uint color;

		void main() {
			gl_Position = vec4(aPos, 0.0, 1.0);
			color = aColor;
		}
	` + "\x00"

	fragmentShaderRed = `
		#version 410
		out vec4 FragColor;

		void main() {
			FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"

	fragmentShaderBlack = `
		#version 410
		out vec4 FragColor;

		void main() {
			FragColor = vec4(0.0, 0.0, 0.0, 1.0);
		}
	` + "\x00"

	fragmentShaderColor = `
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
)

func NewProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func initGlfw(width, height uint) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(int(width), int(height), "game-boy-emulator", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logMessage := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMessage))

		return 0, fmt.Errorf("failed to compile %v: %v", source, logMessage)
	}

	return shader, nil
}
