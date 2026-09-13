package pm

import (
	"errors"
	"math/bits"
)

const FrameSize = 4096
var offsetBits = bits.Len64(FrameSize - 1)

func FrameNumber(addr uint64) uint64 {
	return addr >> offsetBits
}

func FrameOffset(addr uint64) uint64 {
	return addr & (FrameSize - 1)
}

type PhysicalMemory struct {
	data []byte
}

func NewMemory(size uint64) *PhysicalMemory {
	return &PhysicalMemory{
		data: make([]byte, size),
	}
}

func (m *PhysicalMemory) Read(addr uint64) (byte, error) {
	if addr >= uint64(len(m.data)) {
		return 0x0, errors.New("out of range")
	}
	return m.data[addr], nil
}

func (m *PhysicalMemory) Write(addr uint64, value byte) error {
	if addr >= uint64(len(m.data)) {
		return errors.New("out of range")
	}
	m.data[addr] = value
	return nil
}
