package vm

import "errors"

var (
	ErrOutOfRange       = errors.New("out of range")
	ErrPageNotPresented = errors.New("page not presented")
)

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
