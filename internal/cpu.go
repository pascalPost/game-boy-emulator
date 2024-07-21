package internal

import (
	"log"
	"log/slog"
	"unsafe"
)

func isBit7Set(value uint8) bool {
	return (value & 0b1000_0000) != 0
}

func isBit6Set(value uint8) bool {
	return (value & 0b0100_0000) != 0
}

func isBit5Set(value uint8) bool {
	return (value & 0b0010_0000) != 0
}

func isBit4Set(value uint8) bool {
	return (value & 0b0001_0000) != 0
}

func setBit7(value *uint8) {
	*value = *value | 0b1000_0000
}

func setBit6(value *uint8) {
	*value = *value | 0b0100_0000
}

func setBit5(value *uint8) {
	*value = *value | 0b0010_0000
}

func setBit4(value *uint8) {
	*value = *value | 0b0001_0000
}

func clearBit7(value *uint8) {
	*value = *value & 0b0111_1111
}

func clearBit6(value *uint8) {
	*value = *value & 0b1011_1111
}

func clearBit5(value *uint8) {
	*value = *value & 0b1101_1111
}

func clearBit4(value *uint8) {
	*value = *value & 0b1110_1111
}

func clearBits7to4(value *uint8) {
	*value = *value & 0b0000_1111
}

func highPart(value uint16) uint8 {
	return uint8(value >> 8)
}

func lowPart(value uint16) uint8 {
	return uint8(value & 0xff)
}

func highPartPtr(ptr *uint16) *uint8 {
	// this might only work for little endian systems, if so swap with lowPartPtr
	uPtr := unsafe.Pointer(ptr)
	return (*uint8)(unsafe.Pointer(uintptr(uPtr) + 1))
}

func lowPartPtr(ptr *uint16) *uint8 {
	// this might only work for little endian systems, if so swap with highPartPtr
	uPtr := unsafe.Pointer(ptr)
	return (*uint8)(uPtr)
}

type flagsPtr struct {
	data *uint8
}

func (f flagsPtr) clear() {
	clearBits7to4(f.data)
}

func (f flagsPtr) setZ() {
	setBit7(f.data)
}

func (f flagsPtr) setN() {
	setBit6(f.data)
}

func (f flagsPtr) setH() {
	setBit5(f.data)
}

func (f flagsPtr) setC() {
	setBit4(f.data)
}

func (f flagsPtr) z() bool {
	return isBit7Set(*f.data)
}

func (f flagsPtr) n() bool {
	return isBit6Set(*f.data)
}

func (f flagsPtr) h() bool {
	return isBit5Set(*f.data)
}

func (f flagsPtr) c() bool {
	return isBit4Set(*f.data)
}

type registers struct {
	af uint16 // Accumulator & Flags
	bc uint16
	de uint16
	hl uint16
	sp uint16 // Stack Pointer
	pc uint16 // Program Counter/Pointer
	//irie uint16 // Instruction Register & Interrupt Enable
}

func (r *registers) flags() flagsPtr {
	return flagsPtr{lowPartPtr(&r.af)}
}

func (r *registers) aPtr() *uint8 {
	return highPartPtr(&r.af)
}

func (r *registers) bPtr() *uint8 {
	return highPartPtr(&r.bc)
}

func (r *registers) cPtr() *uint8 {
	return lowPartPtr(&r.bc)
}

func (r *registers) dPtr() *uint8 {
	return highPartPtr(&r.de)
}

func (r *registers) ePtr() *uint8 {
	return lowPartPtr(&r.de)
}

func (r *registers) hPtr() *uint8 {
	return highPartPtr(&r.hl)
}

func (r *registers) lPtr() *uint8 {
	return lowPartPtr(&r.hl)
}

func (r *registers) a() uint8 {
	return highPart(r.af)
}

func (r *registers) b() uint8 {
	return highPart(r.bc)
}

func (r *registers) c() uint8 {
	return lowPart(r.bc)
}

func (r *registers) d() uint8 {
	return highPart(r.de)
}

func (r *registers) e() uint8 {
	return lowPart(r.de)
}

func (r *registers) h() uint8 {
	return highPart(r.hl)
}

func (r *registers) l() uint8 {
	return lowPart(r.hl)
}

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

	case 0x01:
		load16BitToRegister(memory, &cpu.registers.pc, &cpu.registers.bc, "BC")
	case 0x11:
		load16BitToRegister(memory, &cpu.registers.pc, &cpu.registers.de, "DE")
	case 0x21:
		load16BitToRegister(memory, &cpu.registers.pc, &cpu.registers.hl, "HL")
	case 0x31:
		load16BitToRegister(memory, &cpu.registers.pc, &cpu.registers.sp, "SP")

	case 0x3E:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.aPtr(), "A")
	case 0x06:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.bPtr(), "B")
	case 0x0E:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.cPtr(), "C")
	case 0x16:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.dPtr(), "D")
	case 0x1E:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.ePtr(), "E")
	case 0x26:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.hPtr(), "H")
	case 0x2E:
		load8BitToRegister(memory, &cpu.registers.pc, cpu.registers.lPtr(), "L")

	case 0x36:
		load8BitToAddressInHL(memory, &cpu.registers.pc, cpu.registers.hl)

	case 0x40:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.b(), "B", "B")
	case 0x41:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.c(), "B", "C")
	case 0x42:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.d(), "B", "D")
	case 0x43:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.e(), "B", "E")
	case 0x44:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.h(), "B", "H")
	case 0x45:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.b(), "B", "L")
	case 0x46:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.hl, "B", "HL")

	case 0x48:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.b(), "C", "B")
	case 0x49:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.c(), "C", "C")
	case 0x4A:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.d(), "C", "D")
	case 0x4B:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.e(), "C", "E")
	case 0x4C:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.h(), "C", "H")
	case 0x4D:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.l(), "C", "L")
	case 0x4E:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.hl, "C", "HL")

	case 0x50:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.b(), "D", "B")
	case 0x51:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.c(), "D", "C")
	case 0x52:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.d(), "D", "D")
	case 0x53:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.e(), "D", "E")
	case 0x54:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.h(), "D", "H")
	case 0x55:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.l(), "D", "L")
	case 0x56:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.hl, "D", "HL")

	case 0x58:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.b(), "E", "B")
	case 0x59:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.c(), "E", "C")
	case 0x5A:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.d(), "E", "D")
	case 0x5B:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.e(), "E", "E")
	case 0x5C:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.h(), "E", "H")
	case 0x5D:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.l(), "E", "L")
	case 0x5E:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.hl, "E", "HL")

	case 0x60:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.b(), "H", "B")
	case 0x61:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.c(), "H", "C")
	case 0x62:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.d(), "H", "D")
	case 0x63:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.e(), "H", "E")
	case 0x64:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.h(), "H", "H")
	case 0x65:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.l(), "H", "L")
	case 0x66:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.hl, "H", "HL")

	case 0x68:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.b(), "L", "B")
	case 0x69:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.c(), "L", "C")
	case 0x6A:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.d(), "L", "D")
	case 0x6B:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.e(), "L", "E")
	case 0x6C:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.h(), "L", "H")
	case 0x6D:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.l(), "L", "L")
	case 0x6E:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.hl, "L", "HL")

	case 0x70:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.bPtr(), cpu.registers.hl, "B")
	case 0x71:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.cPtr(), cpu.registers.hl, "C")
	case 0x72:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.dPtr(), cpu.registers.hl, "D")
	case 0x73:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.ePtr(), cpu.registers.hl, "E")
	case 0x74:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.hPtr(), cpu.registers.hl, "H")
	case 0x75:
		load8BitToAddressInHLFromRegister(memory, cpu.registers.pc, cpu.registers.lPtr(), cpu.registers.hl, "L")

	case 0x7F:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.a(), "A", "A")
	case 0x78:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.b(), "A", "B")
	case 0x79:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.c(), "A", "C")
	case 0x7A:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.d(), "A", "D")
	case 0x7B:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.e(), "A", "E")
	case 0x7C:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.h(), "A", "H")
	case 0x7D:
		load8BitToRegisterFromRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.l(), "A", "L")

	case 0x0A:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.bc, "A", "BC")
	case 0x1A:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.de, "A", "DE")
	case 0x7E:
		load8BitToRegisterFromAddressInRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.hl, "A", "HL")

	case 0xC3:
		jp(memory, &cpu.registers.pc)

	case 0xCD:
		call(memory, &cpu.registers.pc, &cpu.registers.sp)
	case 0xEA:
		loadFromAccumulator(memory, &cpu.registers.pc, cpu.registers.a())

	case 0xB7:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.a(), "A")
	case 0xB0:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.b(), "B")
	case 0xB1:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.c(), "C")
	case 0xB2:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.d(), "D")
	case 0xB3:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.e(), "E")
	case 0xB4:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.h(), "H")
	case 0xB5:
		bitwiseOrRegister(memory, cpu.registers.pc, cpu.registers.aPtr(), cpu.registers.flags(), cpu.registers.l(), "L")

	case 0xC9:
		returnFromFunction(memory, &cpu.registers.pc, &cpu.registers.sp)
	case 0xC0:
		returnFromFunctionConditional(memory, &cpu.registers.pc, &cpu.registers.sp, !cpu.registers.flags().z(), "NZ")
	case 0xC8:
		returnFromFunctionConditional(memory, &cpu.registers.pc, &cpu.registers.sp, cpu.registers.flags().z(), "Z")
	case 0xD0:
		returnFromFunctionConditional(memory, &cpu.registers.pc, &cpu.registers.sp, !cpu.registers.flags().c(), "NC")
	case 0xD8:
		returnFromFunctionConditional(memory, &cpu.registers.pc, &cpu.registers.sp, cpu.registers.flags().c(), "C")

	default:
		log.Panicf("unknown opcode: 0x%02X", opcode)
	}
}
