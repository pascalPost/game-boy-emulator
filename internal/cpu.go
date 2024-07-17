package internal

//type cpu struct {
//	AF [2]byte // Accumulator & Flags, A refers to Hi, Lo contains flags
//	BC [2]byte
//	DE [2]byte
//	HL [2]byte
//	SP [2]byte
//	PC [2]byte
//
// // allocate the memory here (?) would be [0xffff]byte
//
//}
//
//func (c *cpu) A() {
//	return
//}
//
//func readMemory(addr byte) byte {
//
//}
//
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
