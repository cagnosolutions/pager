package pagerv2

import (
	"errors"
)

var (
	ErrIllegalPageAccess = errors.New("illegal page access")
	ErrRecordTooLarge    = errors.New("record too large")
	ErrPageIsFull        = errors.New("page is full")
)

/*
	Unforeseen situations
	---------------------
	page cache is full
*/

// Pagerer is a page cache interface
type Pagerer interface {

	// Read takes the provided page ID and returns the
	// contents of the page. If the page cannot be
	// found in memory, it searches on the disk. If it
	// cannot be found on the disk, an "illegal page
	// access" error is returned. On success, it will
	// return a nil error.
	Read(pid uint32) ([]byte, error)

	// ReadRecord takes the provided page ID and record
	// ID and returns the contents of the record. If the
	// page cannot be found in memory, it searches on the
	// disk. If it cannot be found on the disk an "illegal
	// page access" error is returned. On success, it will
	// return a nil error.
	ReadRecord(pid uint32, rid uint16) ([]byte, error)

	// Write takes data and a page ID and attempts to
	// write the data provided to the page specified
	// using the provided page ID. It returns any
	// potential errors it encounters. On success,
	// it will return a nil error. (Any write made to
	// a page will not be persisted to disk until an
	// explicit call to Sync is made.)
	Write(d []byte, pid uint32) error

	// WriteRecord takes data, a page ID and a record ID
	// and attempts to write the record data to the page
	// specified using the provided page and record IDs.
	// It returns and potential errors it encounters. On
	// success, it will return a nil error. (Any write
	// made to a page or record will not be persisted to
	// disk until an explicit call to Sync is made.)
	WriteRecord(r []byte, pid uint32, rid uint16) error

	// GetFreePageID searches through the pager and returns
	// the first free page ID it finds. If none can be
	// found then dirty pages are flushed to disk to make
	// room, or old pages that are not dirty are evicted
	// to make room. It will always return a free page ID
	// one way or another.
	GetFreePageID() uint32

	// GetFreeRecordID takes a page ID and a required record
	// size and attempts to return a free slot for a record
	// that would fit within the specified page and record
	// size. If a slot cannot be found or if the page is full
	// a "record too large" or "page is full" error will be
	// returned. On success, it will return a nil error.
	GetFreeRecordID(pid uint32, rsize int) (uint16, error)

	// Sync commits the current contents of the pager
	// to stable storage. This means flushing cache's
	// in-memory copy of recently written data to disk.
	Sync() error
}
