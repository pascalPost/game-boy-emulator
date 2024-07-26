package cpu

type Memory struct {
	Data [0xFFFF]byte
}

func (m *Memory) Read(address uint16) uint8 {
	return m.Data[address]
}

func (m *Memory) Write(address uint16, value uint8) {
	m.Data[address] = value
}
