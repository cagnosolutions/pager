package pagerv3

import (
	"errors"
)

// DiskMaxNumPages sets the disk capacity
const DiskMaxNumPages = 15

// DiskManager is responsible for interacting with disk
type DiskManager struct {
	pages map[*PageID]Page
	// tracks the number of pages. -1 indicates that there
	// is no page, and the next to be allocates is 0
	count int
}

// NewDiskManager returns a in-memory mock of disk manager
func NewDiskManager() *DiskManager {
	return &DiskManager{
		pages: make(map[*PageID]Page),
		count: -1,
	}
}

// ReadPage reads a page from pages
func (d *DiskManager) ReadPage(pageID *PageID) (*Page, error) {
	if page, ok := d.pages[pageID]; ok {
		return &page, nil
	}
	return nil, errors.New("Page not found")
}

// WritePage writes a page in memory to pages
func (d *DiskManager) WritePage(p *Page) error {
	d.pages[&p.id] = *p
	return nil
}

// AllocatePage allocates one more page
func (d *DiskManager) AllocatePage() *PageID {
	if d.count == DiskMaxNumPages-1 {
		return nil
	}
	d.count++
	pid := PageID(d.count)
	return &pid
}

// DeallocatePage removes page from disk
func (d *DiskManager) DeallocatePage(pid *PageID) {
	delete(d.pages, pid)
}
