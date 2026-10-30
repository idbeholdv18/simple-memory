package vm

import (
	"errors"
	"github/idbeholdv18/simple-memory/internal/pm"
	"math/bits"
)

const PageSize uint64 = 4096

var (
	ErrOutOfRange        = errors.New("out of range")
	ErrPageNotPresented  = errors.New("page not presented")
	ErrIncorrectVMSize   = errors.New("incorrect vm size")
	ErrPageAlreadyMapped = errors.New("page already mapped")
)

func PageNumber(addr uint64) uint64 {
	return addr >> (bits.Len64(PageSize) - 1)
}

func PageOffset(addr uint64) uint64 {
	return addr & (PageSize - 1)
}

type Page struct {
	isPresented bool
	frame       uint64
}

type PageTable struct {
	table []Page
}

func NewPageTable(pageCount uint64) *PageTable {
	return &PageTable{
		table: make([]Page, pageCount),
	}
}

func (pt *PageTable) Map(page uint64, frame uint64) error {
	if page >= uint64(len(pt.table)) {
		return ErrOutOfRange
	}

	pt.table[page].isPresented = true
	pt.table[page].frame = frame
	return nil
}

func (pt *PageTable) Translate(page uint64) (uint64, error) {
	if page >= uint64(len(pt.table)) {
		return 0, ErrOutOfRange
	}

	if !pt.table[page].isPresented {
		return 0, ErrPageNotPresented
	}

	return pt.table[page].frame, nil
}

func (pt *PageTable) TranslateAddress(addr uint64) (uint64, error) {
	pageNumber := PageNumber(addr)
	pageOffset := PageOffset(addr)

	frameNumber, err := pt.Translate(pageNumber)
	if err != nil {
		return 0, err
	}

	return frameNumber<<uint64(bits.Len64(PageSize)-1) | pageOffset, nil
}

type VirtualMemory struct {
	physicalMemory *pm.PhysicalMemory
	allocator      *pm.FrameAllocator
	pageTable      *PageTable
}

func NewVirtualMemory(pmSize uint64, vmSize uint64) (*VirtualMemory, error) {
	if pmSize == 0 || vmSize == 0 {
		return nil, ErrIncorrectVMSize
	}
	if pmSize%pm.FrameSize != 0 || vmSize%PageSize != 0 {
		return nil, ErrIncorrectVMSize
	}

	frames := pmSize / pm.FrameSize
	pages := vmSize / PageSize

	physicalMemory := pm.NewMemory(pmSize)
	allocator := pm.NewFrameAllocator(frames)
	pageTable := NewPageTable(pages)

	return &VirtualMemory{
		physicalMemory: physicalMemory,
		pageTable:      pageTable,
		allocator:      allocator,
	}, nil
}

func (m *VirtualMemory) MapPage(page uint64) error {
	frame, err := m.pageTable.Translate(page)

	if errors.Is(err, ErrPageNotPresented) {
		if frame, err = m.allocator.Allocate(); err != nil {
			return err
		}

		if err := m.pageTable.Map(page, frame); err != nil {
			if freeErr := m.allocator.Free(frame); freeErr != nil {
				return errors.Join(err, freeErr)
			}
			return err
		}

		return nil
	}

	if err != nil {
		return err
	}
	return ErrPageAlreadyMapped
}
