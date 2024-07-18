package internal

type registers struct {
	a  byte // Accumulator
	b  byte
	c  byte
	d  byte
	e  byte
	f  byte // Flags
	h  byte
	l  byte
	sp [2]byte // Stack Pointer
	pc [2]byte // Program Counter/Pointer
}

func (c *registers) A() byte {
	return c.a
}

//func executeInstruction() {
//	opcode := readMemory(PC)
//}
//
//func writeMemory(addr [2]byte, data byte) {}
//
//func LD_nn_A(PC byte) byte {
//	//  Load from accumulator (direct)
//	// Load to the absolute address specified by the 16-bit operand nn, data from the 8-bit A register.
//	{
//		nn_lsb := readMemory(PC)
//		PC += 1
//	}
//	{
//		nn_msb := readMemory(PC)
//		PC += 1
//	}
//
//	nn := unsigned16Bit(nn_lsb, nn_msb)
//
//	writeMemory(nn, A)
//
//	return PC
//}
