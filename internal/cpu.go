package internal

import (
	"log"
	"log/slog"
)

//type register struct {
//	value uint16
//	name  [2]byte
//}

type registers struct {
	a  uint8 // Accumulator
	b  uint8
	c  uint8
	d  uint8
	e  uint8
	f  uint8 // Flags
	h  uint8
	l  uint8
	sp uint16 // Stack Pointer
	pc uint16 // Program Counter/Pointer
}

//func (r *registers) setBC(value uint16) {
//	lo := uint8(value & 0xff)
//	hi := uint8(value >> 8)
//	r.b = hi
//	r.c = lo
//}

type cpu struct {
	registers registers
}

func (cpu *cpu) runInstruction(memory *memory) {
	opcode := memory.read(cpu.registers.pc)
	slog.Debug("Decode instruction", "PC", fmtHex16(cpu.registers.pc), "Opcode", fmtHex8(opcode))
	cpu.registers.pc++

	// instead of a switch we could also read from an array/slice at the opcode position
	switch opcode {
	case 0x00:
		nop(cpu.registers.pc)
	case 0x31:
		loadTo16BitRegister(memory, &cpu.registers.pc, &cpu.registers.sp, "SP")
	case 0x3E:
		loadTo8BitRegister(memory, &cpu.registers.pc, &cpu.registers.a, "A")
	case 0xC3:
		jp(memory, &cpu.registers.pc)
	case 0xCD:
		call(memory, &cpu.registers.pc, &cpu.registers.sp)
	default:
		log.Panicf("unknown opcode: 0x%02X", opcode)
	}
}
