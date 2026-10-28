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

// Size of per frame
const FrameSize = 4096

// FrameNumber calculates frame number from address
//
// Addr is a full address in physical memory.
// Returns frame number without offset
func FrameNumber(addr uint64) uint64 {
	return addr >> (bits.Len64(FrameSize) - 1)
}

// FrameOffset calculates offset by address inside frame
//
// Addr is a full address in physical memory.
// Returns frame offset without frame number
func FrameOffset(addr uint64) uint64 {
	return addr & (FrameSize - 1)
}

// PhysicalMemory is a plain slice of bytes, representing data
type PhysicalMemory struct {
	data []byte
}

// NewMemory creates a physical memory region of the given size in bytes.
// The size does not have to be a multiple of the frame size.
// Any remaining bytes that do not form a complete frame cannot be
// allocated by the frame allocator.
func NewMemory(size uint64) *PhysicalMemory {
	return &PhysicalMemory{
		data: make([]byte, size),
	}
}

// Read returns byte from physical memory by full address
// Addr have to be in range [0, len(m.data)-1] or it will
// return ErrOutOfRange otherwise.
func (m *PhysicalMemory) Read(addr uint64) (byte, error) {
	if addr >= uint64(len(m.data)) {
		return 0x0, ErrOutOfRange
	}
	return m.data[addr], nil
}

// Write stores a byte at the given physical address.
// The address must be in the range [0, len(m.data)-1].
//
// It returns ErrOutOfRange if the address is outside this range.
//
// It returns nil on successful write.
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
