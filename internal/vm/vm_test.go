package vm

import (
	"errors"
	"testing"
)

func TestNewPageTable(t *testing.T) {
	pt := NewPageTable(3)
	if len(pt.table) != 3 {
		t.Fatalf("incorrect page table size: %d, expected 3", len(pt.table))
	}
}

func TestMap(t *testing.T) {
	pt := NewPageTable(3)

	err := pt.Map(0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pt.table[0].frame != 10 || !pt.table[0].isPresented {
		t.Errorf("expected 10:true, got %d:%t", pt.table[0].frame, pt.table[0].isPresented)
	}

	err = pt.Map(0, 2)
	if pt.table[0].frame != 2 || !pt.table[0].isPresented {
		t.Errorf("expected 2:true, got %d:%t", pt.table[0].frame, pt.table[0].isPresented)
	}
}

func TestTranslate(t *testing.T) {
	pt := NewPageTable(3)

	_, err := pt.Translate(0)
	if !errors.Is(err, ErrPageNotPresented) {
		t.Fatalf("expected err: %v, got %v", ErrPageNotPresented, err)
	}

	pt.Map(0, 10)

	frame, err := pt.Translate(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if frame != 10 {
		t.Errorf("expected 10, got %d", frame)
	}
}

func TestTranslateAddress(t *testing.T) {
	var physicalAddress uint64 = 0x4321

	pt := NewPageTable(4)
	pt.Map(0, 4)

	res, err := pt.TranslateAddress(0x0321)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res != physicalAddress {
		t.Errorf("expected %x, got %x", physicalAddress, res)
	}
}
