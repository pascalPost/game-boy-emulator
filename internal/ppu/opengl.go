package ppu

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
	"unsafe"
)

func ActivateDebugOutput() {
	var flags int32
	gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
	if flags&gl.CONTEXT_FLAG_DEBUG_BIT != 0 {
		log.Println("OpenGL debug context available")

		debugCallback := func(
			source uint32,
			gltype uint32,
			id uint32,
			severity uint32,
			length int32,
			message string,
			userParam unsafe.Pointer) {

			// ignore non-significant error/warning codes
			if id == 131169 || id == 131185 || id == 131218 || id == 131204 {
				return
			}

			println("Debug message (source): ", source)
			println("Debug message (type): ", gltype)
			println("Debug message (id): ", id)
			println("Debug message (severity): ", severity)
			println("Debug message (length): ", length)
			println("Debug message: ", message)
		}

		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
		gl.DebugMessageCallback(debugCallback, nil)
		gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
	}

}
