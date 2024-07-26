package cpu

type Registers struct {
	AF uint16 // Accumulator & Flags
	BC uint16
	DE uint16
	HL uint16
	SP uint16 // Stack Pointer
	PC uint16 // Program Counter/Pointer
	//irie uint16 // Instruction Register & Interrupt Enable
}

func (r *Registers) Flags() FlagsPtr {
	return NewFlagsPtr(LowPartPtr(&r.AF))
}

func (r *Registers) APtr() *uint8 {
	return HighPartPtr(&r.AF)
}

func (r *Registers) BPtr() *uint8 {
	return HighPartPtr(&r.BC)
}

func (r *Registers) CPtr() *uint8 {
	return LowPartPtr(&r.BC)
}

func (r *Registers) DPtr() *uint8 {
	return HighPartPtr(&r.DE)
}

func (r *Registers) EPtr() *uint8 {
	return LowPartPtr(&r.DE)
}

func (r *Registers) HPtr() *uint8 {
	return HighPartPtr(&r.HL)
}

func (r *Registers) LPtr() *uint8 {
	return LowPartPtr(&r.HL)
}

func (r *Registers) A() uint8 {
	return HighPart(r.AF)
}

func (r *Registers) B() uint8 {
	return HighPart(r.BC)
}

func (r *Registers) C() uint8 {
	return LowPart(r.BC)
}

func (r *Registers) D() uint8 {
	return HighPart(r.DE)
}

func (r *Registers) E() uint8 {
	return LowPart(r.DE)
}

func (r *Registers) H() uint8 {
	return HighPart(r.HL)
}

func (r *Registers) L() uint8 {
	return LowPart(r.HL)
}

type Cpu struct {
	Registers Registers
	IME       bool // IME (interrupt master enable) flag indicating if interrupts are enabled (1) or disabled (0)
}
