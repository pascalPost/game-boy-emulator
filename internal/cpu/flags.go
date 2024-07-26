package cpu

type FlagsPtr struct {
	data *uint8
}

func NewFlagsPtr(data *uint8) FlagsPtr {
	return FlagsPtr{data}
}

func (f FlagsPtr) ClearAll() {
	clearBits7to4(f.data)
}

func (f FlagsPtr) ClearZ() {
	_ = clearBit(f.data, 7)
}

func (f FlagsPtr) ClearN() {
	_ = clearBit(f.data, 6)
}

func (f FlagsPtr) ClearH() {
	_ = clearBit(f.data, 5)
}

func (f FlagsPtr) clearC() {
	_ = clearBit(f.data, 4)
}

func (f FlagsPtr) SetZ() {
	_ = setBit(f.data, 7)
}

func (f FlagsPtr) SetN() {
	_ = setBit(f.data, 6)
}

func (f FlagsPtr) SetH() {
	_ = setBit(f.data, 5)
}

func (f FlagsPtr) SetC() {
	_ = setBit(f.data, 4)
}

func (f FlagsPtr) Z() bool {
	return checkBit(*f.data, 7)
}

func (f FlagsPtr) n() bool {
	return checkBit(*f.data, 6)
}

func (f FlagsPtr) h() bool {
	return checkBit(*f.data, 5)
}

func (f FlagsPtr) C() bool {
	return checkBit(*f.data, 4)
}
