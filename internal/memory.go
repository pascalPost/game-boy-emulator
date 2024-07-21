package internal

type memory struct {
	data [0xFFFF]byte
}

func (m *memory) read(address uint16) uint8 {
	return m.data[address]
}

func (m *memory) write(address uint16, value uint8) {
	m.data[address] = value
}
