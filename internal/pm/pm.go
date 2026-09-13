package pm

import (
	"errors"
	"math/bits"
)

var (
	ErrOutOfRange       = errors.New("out of range")
	ErrOutOfMemory      = errors.New("out of memory")
	ErrFreeFrameRelease = errors.New("free frame release")
)

const FrameSize = 4096

func FrameNumber(addr uint64) uint64 {
	return addr >> bits.Len64(FrameSize-1)
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
		return 0x0, ErrOutOfRange
	}
	return m.data[addr], nil
}

func (m *PhysicalMemory) Write(addr uint64, value byte) error {
	if addr >= uint64(len(m.data)) {
		return ErrOutOfRange
	}
	m.data[addr] = value
	return nil
}

type FrameAllocator struct {
	state []bool
}

func NewFrameAllocator(frameCount uint64) *FrameAllocator {
	return &FrameAllocator{
		state: make([]bool, frameCount),
	}
}

func (a *FrameAllocator) Allocate() (uint64, error) {
	for i, isBusy := range a.state {
		if !isBusy {
			a.state[i] = true
			return uint64(i), nil
		}
	}
	return 0, ErrOutOfMemory
}

func (a *FrameAllocator) Free(frame uint64) error {
	if frame >= uint64(len(a.state)) {
		return ErrOutOfRange
	}

	if !a.state[frame] {
		return ErrFreeFrameRelease
	}

	a.state[frame] = false
	return nil
}
