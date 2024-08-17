package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pascalPost/game-boy-emulator/internal/ppu"
	"runtime"
)

const (
	vertexShader2DTexture = `
		#version 410 core
	
		layout(location = 0) in vec2 aPos;
		layout(location = 1) in vec2 aTexCoord;
	
		out vec2 TexCoord;
	
		void main()
		{
		    TexCoord = aTexCoord;
		    gl_Position = vec4(aPos, 0.0, 1.0);
		}
	` + "\x00"

	vertexShader = `
		#version 410 core
		
		layout(location = 0) in vec2 vPos;
		
		void main()
		{
		    gl_Position = vec4(vPos, 0.0, 1.0);
		}
	` + "\x00"

	fragmentShader = `
		#version 410 core
		out vec4 FragColor;
		
		void main()
		{
			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		}
	` + "\x00"

	fragmentShaderTexture = `
		#version 410 core
	
		in vec2 TexCoord;
		out vec4 FragColor;
	
		uniform sampler2D screenTexture;
	
		void main()
		{
		    FragColor = texture(screenTexture, TexCoord);
		}
	` + "\x00"
)

func main() {
	runtime.LockOSThread()

	ppu.InitGlfw()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	window, err := glfw.CreateWindow(int(500), int(500), "window1", nil, nil)
	if err != nil {
		panic(err)
	}

	window2, err := glfw.CreateWindow(int(500), int(500), "window2", nil, window)
	if err != nil {
		panic(err)
	}

	ppu.InitOpenGL()

	window.MakeContextCurrent()
	ppu.ActivateDebugOutput()

	window2.MakeContextCurrent()
	ppu.ActivateDebugOutput()

	window.MakeContextCurrent()

	program := ppu.NewProgram(vertexShader, fragmentShader)

	textureProgram := ppu.NewProgram(vertexShader2DTexture, fragmentShaderTexture)

	vposLocation := gl.GetAttribLocation(program, gl.Str("vPos\x00"))

	vertexData := []float32{
		-0.5, -0.5,
		0.0, 0.5,
		0.5, -0.5,
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	defer func() {
		gl.DeleteVertexArrays(1, &vao)
	}()

	gl.BindVertexArray(vao)

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, ppu.Float32ByteSize*len(vertexData), gl.Ptr(vertexData), gl.STATIC_DRAW)

	gl.VertexAttribPointer(uint32(vposLocation), 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(uint32(vposLocation))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// Create framebuffer object (FBO)
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	defer func() {
		gl.DeleteFramebuffers(1, &fbo)
	}()

	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

	// Create texture to render to
	var texColorBuffer uint32
	gl.GenTextures(1, &texColorBuffer)
	gl.BindTexture(gl.TEXTURE_2D, texColorBuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 500, 500, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texColorBuffer, 0)

	//var rbo uint32
	//gl.GenRenderbuffers(1, &rbo)
	//gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	//gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, 500, 500)
	//gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("Framebuffer is not complete")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// quad that fills the entire screen in normalized device coordinates
	quadVertices := []float32{
		// positions   // texCoords
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 0.0,

		-1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, 1.0, 1.0,
	}

	var quadVAO, quadVBO uint32
	gl.GenVertexArrays(1, &quadVAO)
	defer func() {
		gl.DeleteVertexArrays(1, &quadVAO)
	}()

	gl.GenBuffers(1, &quadVBO)
	gl.BindVertexArray(quadVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*ppu.Float32ByteSize, gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	window2.MakeContextCurrent()

	var quadVAO2 uint32
	gl.GenVertexArrays(1, &quadVAO2)
	defer func() {
		gl.DeleteVertexArrays(1, &quadVAO2)
	}()

	gl.BindVertexArray(quadVAO2)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	gl.BindVertexArray(0)

	window.MakeContextCurrent()

	// place second window to the right of the first window
	{
		x, y := window.GetPos()
		width, _ := window.GetSize()

		window2.SetPos(x+width, y)
	}

	window2.MakeContextCurrent()

	gl.UseProgram(program)

	var vao2 uint32
	gl.GenVertexArrays(1, &vao2)
	defer func() {
		gl.DeleteVertexArrays(1, &vao2)
	}()

	gl.BindVertexArray(vao2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.EnableVertexAttribArray(uint32(vposLocation))
	gl.VertexAttribPointer(uint32(vposLocation), 2, gl.FLOAT, false, 0, nil)

	for !window.ShouldClose() && !window2.ShouldClose() {
		//Render to the framebuffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
		gl.Viewport(0, 0, 500, 500)
		gl.ClearColor(1.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexData)/ppu.Dimensions))
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0) // unbind fbo

		// Render to the first window

		{
			width, height := window.GetFramebufferSize()
			window.MakeContextCurrent()
			gl.Viewport(0, 0, int32(width), int32(height))
			gl.ClearColor(0.0, 1.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.UseProgram(textureProgram)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, texColorBuffer)
			gl.BindVertexArray(quadVAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)

			window.SwapBuffers()
			println("End window1: ", gl.GetError())
		}

		//Render the FBO texture to the second window with zoom
		{
			width, height := window2.GetFramebufferSize()
			window2.MakeContextCurrent()
			gl.Viewport(0, 0, int32(width), int32(height))
			gl.ClearColor(1.0, 1.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.UseProgram(textureProgram)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, texColorBuffer)
			gl.BindVertexArray(quadVAO2)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)

			//gl.UseProgram(program)
			//gl.BindVertexArray(vao2)
			//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertexData)/ppu.Dimensions))

			window2.SwapBuffers()
			println("End window2: ", gl.GetError())
		}

		//println(gl.GetError())

		glfw.PollEvents()
	}

	//gl.DeleteTextures(1, &texColorBuffer)
}
