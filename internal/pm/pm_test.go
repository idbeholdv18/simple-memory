package pm

import (
	"errors"
	"testing"
)

func TestNewMemory(t *testing.T) {
	mem := NewMemory(1024)

	if len(mem.data) != 1024 {
		t.Errorf("expected memory size is 1024, got %d", len(mem.data))
	}
}

func TestWriteAndRead(t *testing.T) {
	mem := NewMemory(1024)

	if err := mem.Write(10, 12); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	value, err := mem.Read(10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != 12 {
		t.Errorf("Expected value is 12, got %d", value)
	}
}

func TestReadOutOfRange(t *testing.T) {
	mem := NewMemory(1024)

	_, err := mem.Read(1024)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestWriteOutOfRange(t *testing.T) {
	mem := NewMemory(1024)

	err := mem.Write(1024, 10)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFrameNumbersAndOffsets(t *testing.T) {
	var tests = []struct {
		name string
		addr uint64
		want [2]uint64
	}{
		{name: "0 should be 0:0", addr: 0, want: [2]uint64{0, 0}},
		{name: "1 should be 0:1", addr: 1, want: [2]uint64{0, 1}},
		{name: "4095 should be 0:4095", addr: 4095, want: [2]uint64{0, 4095}},
		{name: "4096 should be 1:0", addr: 4096, want: [2]uint64{1, 0}},
		{name: "4097 should be 1:1", addr: 4097, want: [2]uint64{1, 1}},
		{name: "8191 should be 1:4095", addr: 8191, want: [2]uint64{1, 4095}},
		{name: "8192 should be 2:0", addr: 8192, want: [2]uint64{2, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frameNumber := FrameNumber(tt.addr)
			frameOffset := FrameOffset(tt.addr)
			if frameNumber != tt.want[0] || frameOffset != tt.want[1] {
				t.Errorf("expected %d:%d, got %d:%d", tt.want[0], tt.want[1], frameNumber, frameOffset)
			}
		})
	}
}

func TestNewFrameAllocator(t *testing.T) {
	a := NewFrameAllocator(3)

	if len(a.state) != 3 {
		t.Errorf("incorrect allocator state size: %d, expected 3", len(a.state))
	}
}

func TestAllocate(t *testing.T) {
	var frames uint64 = 3
	a := NewFrameAllocator(frames)

	var err error

	for i := range frames {
		frame, err := a.Allocate()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if frame != i {
			t.Errorf("incorrect frame was allocated: %d, expected %d", frame, i)
		}
	}

	_, err = a.Allocate()
	if !errors.Is(err, ErrOutOfMemory) {
		t.Errorf("expected err: %v", ErrOutOfMemory)
	}
}

func TestFree(t *testing.T) {
	a := NewFrameAllocator(1)

	if err := a.Free(0); !errors.Is(err, ErrFreeFrameRelease) {
		t.Errorf("expected error: %v, got %v", ErrFreeFrameRelease, err)
	}

	frame, err := a.Allocate()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = a.Free(frame)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = a.Free(frame)
	if !errors.Is(err, ErrFreeFrameRelease) {
		t.Errorf("expected err: %v", ErrFreeFrameRelease)
	}
}

func TestFirstFitAllocate(t *testing.T) {
	var frames uint64 = 3
	a := NewFrameAllocator(frames)

	var err error

	for range frames {
		_, err = a.Allocate()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	a.Free(1)
	frame, _ := a.Allocate()
	if frame != 1 {
		t.Errorf("incorrect frame was allocated: %d, expected 1", frame)
	}
}
