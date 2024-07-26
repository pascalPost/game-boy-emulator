package instructions

import (
	"fmt"
	. "github.com/pascalPost/game-boy-emulator/internal/cpu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions/alu"
	"github.com/pascalPost/game-boy-emulator/internal/cpu/instructions/rotate"
	"log"
	"log/slog"
)

func RunInstruction(cpu *Cpu, memory *Memory) {
	opcode := memory.Read(cpu.Registers.PC)
	if opcode != uint8(0xCB) {
		slog.Debug("Decode instruction", "PC", FmtHex16(cpu.Registers.PC), "Opcode", FmtHex8(opcode))
	}
	cpu.Registers.PC++

	// instead of a switch we could also read from an array/slice at the opcode position
	switch opcode {
	case 0xCB:
		runPrefixedInstruction(cpu, memory)

	case 0x00:
		Nop(cpu.Registers.PC)

	case 0x01:
		Load16BitToRegister(memory, &cpu.Registers.PC, &cpu.Registers.BC, "BC")
	case 0x11:
		Load16BitToRegister(memory, &cpu.Registers.PC, &cpu.Registers.DE, "DE")
	case 0x21:
		Load16BitToRegister(memory, &cpu.Registers.PC, &cpu.Registers.HL, "HL")
	case 0x31:
		Load16BitToRegister(memory, &cpu.Registers.PC, &cpu.Registers.SP, "SP")

	case 0x3E:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.APtr(), "A")
	case 0x06:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.BPtr(), "B")
	case 0x0E:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.CPtr(), "C")
	case 0x16:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.DPtr(), "D")
	case 0x1E:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.EPtr(), "E")
	case 0x26:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.HPtr(), "H")
	case 0x2E:
		Load8BitToRegister(memory, &cpu.Registers.PC, cpu.Registers.LPtr(), "L")

	case 0x36:
		Load8BitToAddressInHL(memory, &cpu.Registers.PC, cpu.Registers.HL)

	case 0x40:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.B(), "B", "B")
	case 0x41:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.C(), "B", "C")
	case 0x42:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.D(), "B", "D")
	case 0x43:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.E(), "B", "E")
	case 0x44:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.H(), "B", "H")
	case 0x45:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.B(), "B", "L")
	case 0x46:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.HL, "B", "HL")
	case 0x47:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.A(), "B", "A")

	case 0x48:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.B(), "C", "B")
	case 0x49:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.C(), "C", "C")
	case 0x4A:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.D(), "C", "D")
	case 0x4B:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.E(), "C", "E")
	case 0x4C:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.H(), "C", "H")
	case 0x4D:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.L(), "C", "L")
	case 0x4E:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.HL, "C", "HL")
	case 0x4F:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.A(), "C", "A")

	case 0x50:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.B(), "D", "B")
	case 0x51:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.C(), "D", "C")
	case 0x52:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.D(), "D", "D")
	case 0x53:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.E(), "D", "E")
	case 0x54:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.H(), "D", "H")
	case 0x55:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.L(), "D", "L")
	case 0x56:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.HL, "D", "HL")
	case 0x57:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.A(), "D", "A")

	case 0x58:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.B(), "E", "B")
	case 0x59:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.C(), "E", "C")
	case 0x5A:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.D(), "E", "D")
	case 0x5B:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.E(), "E", "E")
	case 0x5C:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.H(), "E", "H")
	case 0x5D:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.L(), "E", "L")
	case 0x5E:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.HL, "E", "HL")
	case 0x5F:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.A(), "E", "A")

	case 0x60:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.B(), "H", "B")
	case 0x61:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.C(), "H", "C")
	case 0x62:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.D(), "H", "D")
	case 0x63:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.E(), "H", "E")
	case 0x64:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.H(), "H", "H")
	case 0x65:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.L(), "H", "L")
	case 0x66:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.HL, "H", "HL")
	case 0x67:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.A(), "H", "A")

	case 0x68:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.B(), "L", "B")
	case 0x69:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.C(), "L", "C")
	case 0x6A:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.D(), "L", "D")
	case 0x6B:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.E(), "L", "E")
	case 0x6C:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.H(), "L", "H")
	case 0x6D:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.L(), "L", "L")
	case 0x6E:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.HL, "L", "HL")
	case 0x6F:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.A(), "L", "A")

	case 0x70:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.HL, "B")
	case 0x71:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.HL, "C")
	case 0x72:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.HL, "D")
	case 0x73:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.HL, "E")
	case 0x74:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.HL, "H")
	case 0x75:
		Load8BitToAddressInHLFromRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.HL, "L")

	case 0x7F:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.A(), "A", "A")

	case 0x78:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.B(), "A", "B")
	case 0x79:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.C(), "A", "C")
	case 0x7A:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.D(), "A", "D")
	case 0x7B:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.E(), "A", "E")
	case 0x7C:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.H(), "A", "H")
	case 0x7D:
		Load8BitToRegisterFromRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.L(), "A", "L")

	case 0x0A:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.BC, "A", "BC")
	case 0x1A:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.DE, "A", "DE")
	case 0x7E:
		Load8BitToRegisterFromAddressInRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.HL, "A", "HL")

	case 0xC3:
		Jp(memory, &cpu.Registers.PC)

	case 0x18:
		RelativeJump(memory, &cpu.Registers.PC)

	case 0x20:
		RelativeJumpConditional(memory, &cpu.Registers.PC, !cpu.Registers.Flags().Z(), "NZ")
	case 0x28:
		RelativeJumpConditional(memory, &cpu.Registers.PC, cpu.Registers.Flags().Z(), "Z")
	case 0x30:
		RelativeJumpConditional(memory, &cpu.Registers.PC, !cpu.Registers.Flags().C(), "NC")
	case 0x38:
		RelativeJumpConditional(memory, &cpu.Registers.PC, cpu.Registers.Flags().C(), "c")

	case 0x02:
		LoadFromRegisterIndirect(memory, cpu.Registers.PC, cpu.Registers.BC, cpu.Registers.A(), "BC", "A")
	case 0x12:
		LoadFromRegisterIndirect(memory, cpu.Registers.PC, cpu.Registers.DE, cpu.Registers.A(), "DE", "A")
	case 0x77:
		LoadFromRegisterIndirect(memory, cpu.Registers.PC, cpu.Registers.HL, cpu.Registers.A(), "HL", "A")
	case 0xE0:
		LoadFromAccumulatorDirectH(memory, &cpu.Registers.PC, cpu.Registers.A())
	case 0xEA:
		LoadFromAccumulatorDirect(memory, &cpu.Registers.PC, cpu.Registers.A())

	case 0xE2:
		LoadFromAccumulatorIndirectHC(memory, cpu.Registers.PC, cpu.Registers.A(), cpu.Registers.C())

	case 0x2A:
		LoadAccumulatorIndirectHLIncrement(memory, cpu.Registers.PC, cpu.Registers.APtr(), &cpu.Registers.HL)

	case 0x32:
		LoadFromAccumulatorIndirectHLDecrement(memory, cpu.Registers.PC, cpu.Registers.A(), &cpu.Registers.HL)

	case 0x3A:
		LoadAccumulatorIndirectHLDecrement(memory, cpu.Registers.PC, cpu.Registers.APtr(), &cpu.Registers.HL)

	case 0xF0:
		LoadAccumulatorDirectLeastSignificantByte(memory, &cpu.Registers.PC, cpu.Registers.APtr())

	case 0xB7:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.A(), "A")
	case 0xB0:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.B(), "B")
	case 0xB1:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.C(), "C")
	case 0xB2:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.D(), "D")
	case 0xB3:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.E(), "E")
	case 0xB4:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.H(), "H")
	case 0xB5:
		alu.BitwiseOrRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.L(), "L")

	case 0xAF:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.A(), "A")
	case 0xA8:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.B(), "B")
	case 0xA9:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.C(), "C")
	case 0xAA:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.D(), "D")
	case 0xAB:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.E(), "E")
	case 0xAC:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.H(), "H")
	case 0xAD:
		alu.BitwiseXorRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), cpu.Registers.L(), "L")

	case 0xC9:
		ReturnFromFunction(memory, &cpu.Registers.PC, &cpu.Registers.SP)
	case 0xC0:
		ReturnFromFunctionConditional(memory, &cpu.Registers.PC, &cpu.Registers.SP, !cpu.Registers.Flags().Z(), "NZ")
	case 0xC8:
		ReturnFromFunctionConditional(memory, &cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.Flags().Z(), "Z")
	case 0xD0:
		ReturnFromFunctionConditional(memory, &cpu.Registers.PC, &cpu.Registers.SP, !cpu.Registers.Flags().C(), "NC")
	case 0xD8:
		ReturnFromFunctionConditional(memory, &cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.Flags().C(), "C")

	case 0x87:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.A(), cpu.Registers.Flags(), "A")
	case 0x80:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.B(), cpu.Registers.Flags(), "B")
	case 0x81:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.C(), cpu.Registers.Flags(), "C")
	case 0x82:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.D(), cpu.Registers.Flags(), "D")
	case 0x83:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.E(), cpu.Registers.Flags(), "E")
	case 0x84:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.H(), cpu.Registers.Flags(), "H")
	case 0x85:
		alu.AddRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.L(), cpu.Registers.Flags(), "L")
	case 0x86:
		alu.AddIndirectHL(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.HL, cpu.Registers.Flags())
	case 0xC6:
		alu.AddImmediate(memory, &cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())

	case 0x97:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.A(), cpu.Registers.Flags(), "A")
	case 0x90:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.B(), cpu.Registers.Flags(), "B")
	case 0x91:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.C(), cpu.Registers.Flags(), "C")
	case 0x92:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.D(), cpu.Registers.Flags(), "D")
	case 0x93:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.E(), cpu.Registers.Flags(), "E")
	case 0x94:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.H(), cpu.Registers.Flags(), "H")
	case 0x95:
		alu.SubtractRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.L(), cpu.Registers.Flags(), "L")
	case 0x96:
		alu.SubtractIndirectHL(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.HL, cpu.Registers.Flags())
	case 0xD6:
		alu.SubtractImmediate(memory, &cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())

	case 0x3C:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), "A")
	case 0x04:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.Flags(), "B")
	case 0x0C:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.Flags(), "C")
	case 0x14:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.Flags(), "D")
	case 0x1C:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.Flags(), "E")
	case 0x24:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.Flags(), "H")
	case 0x2C:
		alu.IncrementRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.Flags(), "L")

	case 0x03:
		alu.Increment16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.BC, "BC")
	case 0x13:
		alu.Increment16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.DE, "DE")
	case 0x23:
		alu.Increment16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.HL, "HL")
	case 0x33:
		alu.Increment16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.SP, "SP")

	case 0x0B:
		alu.Decrement16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.BC, "BC")
	case 0x1B:
		alu.Decrement16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.DE, "DE")
	case 0x2B:
		alu.Decrement16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.HL, "HL")
	case 0x3B:
		alu.Decrement16BitRegister(memory, cpu.Registers.PC, &cpu.Registers.SP, "SP")

	case 0xFE:
		alu.CompareImmediate(memory, &cpu.Registers.PC, cpu.Registers.A(), cpu.Registers.Flags())

	case 0xCD:
		Call(memory, &cpu.Registers.PC, &cpu.Registers.SP)

	case 0xF5:
		PushToStack(memory, cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.AF, "AF")
	case 0xC5:
		PushToStack(memory, cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.BC, "BC")
	case 0xD5:
		PushToStack(memory, cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.DE, "DE")
	case 0xE5:
		PushToStack(memory, cpu.Registers.PC, &cpu.Registers.SP, cpu.Registers.HL, "HL")

	case 0xF1:
		PopFromStack(memory, cpu.Registers.PC, &cpu.Registers.SP, &cpu.Registers.AF, "AF")
	case 0xC1:
		PopFromStack(memory, cpu.Registers.PC, &cpu.Registers.SP, &cpu.Registers.BC, "BC")
	case 0xD1:
		PopFromStack(memory, cpu.Registers.PC, &cpu.Registers.SP, &cpu.Registers.DE, "DE")
	case 0xE1:
		PopFromStack(memory, cpu.Registers.PC, &cpu.Registers.SP, &cpu.Registers.HL, "HL")

	case 0xF3:
		DisableInterrupts(memory, cpu.Registers.PC, &cpu.IME)

	case 0x07:
		rotate.LeftCircularAccumulator(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())
	case 0x17:
		rotate.LeftAccumulator(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())
	case 0x0F:
		rotate.RightCircularAccumulator(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())
	case 0x1F:
		rotate.RightAccumulator(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags())

	default:
		log.Panicf("unknown opcode: 0x%02X", opcode)
	}
}

func runPrefixedInstruction(cpu *Cpu, memory *Memory) {
	opcode := memory.Read(cpu.Registers.PC)
	slog.Debug("Decode prefixed instruction", "PC", FmtHex16(cpu.Registers.PC), "Opcode", fmt.Sprintf("0x%02X 0x%02X", 0xCB, opcode))
	cpu.Registers.PC++

	switch opcode {
	case 0x40:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.B(), "B")
	case 0x41:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.C(), "C")
	case 0x42:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.D(), "D")
	case 0x43:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.E(), "E")
	case 0x44:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.H(), "H")
	case 0x45:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.L(), "L")
	case 0x47:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 0, cpu.Registers.A(), "A")

	case 0x48:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.B(), "B")
	case 0x49:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.C(), "C")
	case 0x4A:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.D(), "D")
	case 0x4B:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.E(), "E")
	case 0x4C:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.H(), "H")
	case 0x4D:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.L(), "L")
	case 0x4F:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 1, cpu.Registers.A(), "A")

	case 0x50:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.B(), "B")
	case 0x51:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.C(), "C")
	case 0x52:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.D(), "D")
	case 0x53:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.E(), "E")
	case 0x54:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.H(), "H")
	case 0x55:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.L(), "L")
	case 0x57:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 2, cpu.Registers.A(), "A")

	case 0x58:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.B(), "B")
	case 0x59:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.C(), "C")
	case 0x5A:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.D(), "D")
	case 0x5B:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.E(), "E")
	case 0x5C:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.H(), "H")
	case 0x5D:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.L(), "L")
	case 0x5F:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 3, cpu.Registers.A(), "A")

	case 0x60:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.B(), "B")
	case 0x61:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.C(), "C")
	case 0x62:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.D(), "D")
	case 0x63:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.E(), "E")
	case 0x64:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.H(), "H")
	case 0x65:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.L(), "L")
	case 0x67:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 4, cpu.Registers.A(), "A")

	case 0x68:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.B(), "B")
	case 0x69:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.C(), "C")
	case 0x6A:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.D(), "D")
	case 0x6B:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.E(), "E")
	case 0x6C:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.H(), "H")
	case 0x6D:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.L(), "L")
	case 0x6F:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 5, cpu.Registers.A(), "A")

	case 0x70:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.B(), "B")
	case 0x71:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.C(), "C")
	case 0x72:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.D(), "D")
	case 0x73:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.E(), "E")
	case 0x74:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.H(), "H")
	case 0x75:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.L(), "L")
	case 0x77:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 6, cpu.Registers.A(), "A")

	case 0x78:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.B(), "B")
	case 0x79:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.C(), "C")
	case 0x7A:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.D(), "D")
	case 0x7B:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.E(), "E")
	case 0x7C:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.H(), "H")
	case 0x7D:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.L(), "L")
	case 0x7F:
		TestBitRegister(memory, cpu.Registers.PC, cpu.Registers.Flags(), 7, cpu.Registers.A(), "A")

	case 0x07:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), "A")
	case 0x00:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.Flags(), "B")
	case 0x01:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.Flags(), "C")
	case 0x02:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.Flags(), "D")
	case 0x03:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.Flags(), "E")
	case 0x04:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.Flags(), "H")
	case 0x05:
		rotate.LeftCircularRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.Flags(), "L")
	case 0x06:
		rotate.LeftCircularIndirectHL(memory, cpu.Registers.PC, cpu.Registers.HL, cpu.Registers.Flags())

	case 0x17:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), "A")
	case 0x10:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.Flags(), "B")
	case 0x11:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.Flags(), "C")
	case 0x12:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.Flags(), "D")
	case 0x13:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.Flags(), "E")
	case 0x14:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.Flags(), "H")
	case 0x15:
		rotate.LeftRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.Flags(), "L")
	case 0x16:
		rotate.LeftIndirectHL(memory, cpu.Registers.PC, cpu.Registers.HL, cpu.Registers.Flags())

	case 0x0F:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), "A")
	case 0x08:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.Flags(), "B")
	case 0x09:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.Flags(), "C")
	case 0x0A:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.Flags(), "D")
	case 0x0B:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.Flags(), "E")
	case 0x0C:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.Flags(), "H")
	case 0x0D:
		rotate.RightCircularRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.Flags(), "L")
	case 0x0E:
		rotate.RightCircularIndirectHL(memory, cpu.Registers.PC, cpu.Registers.HL, cpu.Registers.Flags())

	case 0x1F:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.APtr(), cpu.Registers.Flags(), "A")
	case 0x18:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.BPtr(), cpu.Registers.Flags(), "B")
	case 0x19:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.CPtr(), cpu.Registers.Flags(), "C")
	case 0x1A:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.DPtr(), cpu.Registers.Flags(), "D")
	case 0x1B:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.EPtr(), cpu.Registers.Flags(), "E")
	case 0x1C:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.HPtr(), cpu.Registers.Flags(), "H")
	case 0x1D:
		rotate.RightRegister(memory, cpu.Registers.PC, cpu.Registers.LPtr(), cpu.Registers.Flags(), "L")
	case 0x1E:
		rotate.RightIndirectHL(memory, cpu.Registers.PC, cpu.Registers.HL, cpu.Registers.Flags())

	default:
		log.Panicf("unknown prefixed opcode: 0x%02X", opcode)
	}
}
