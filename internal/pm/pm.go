package pm

import (
	"errors"
	"math/bits"
)

var (
	ERR_OUT_OF_RANGE       = errors.New("out of range")
	ERR_OUT_OF_MEMORY      = errors.New("out of memory")
	ERR_FREE_FRAME_RELEASE = errors.New("free frame release")
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
		return 0x0, ERR_OUT_OF_RANGE
	}
	return m.data[addr], nil
}

func (m *PhysicalMemory) Write(addr uint64, value byte) error {
	if addr >= uint64(len(m.data)) {
		return ERR_OUT_OF_RANGE
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
	return 0, ERR_OUT_OF_MEMORY
}

func (a *FrameAllocator) Free(frame uint64) error {
	if frame >= uint64(len(a.state)) {
		return ERR_OUT_OF_RANGE
	}

	if !a.state[frame] {
		return ERR_FREE_FRAME_RELEASE
	}

	a.state[frame] = false
	return nil
}
