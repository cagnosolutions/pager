package pager

import (
	"io"
	"os"
	"path/filepath"
)

// PageManager is a slotted page PageManager manager
type PageManager struct {
	name        string
	fp          *os.File
	pageHeaders []*pageHeader
	freePages   int
	pids        *autoPageID
}

// OpenPageManager opens an existing PageManager at the location
// provided, or creates and returns a new PageManager at
// the path provided.
func OpenPageManager(path string) (*PageManager, error) {
	// sanitize path
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// split path
	dir, name := filepath.Split(filepath.ToSlash(path))
	// init PageManager and dirs
	var fp *os.File
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
		// create PageManager
		fp, err = os.Create(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		// close PageManager
		err = fp.Close()
		if err != nil {
			return nil, err
		}
	}
	// open existing PageManager
	fp, err = os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	// create page PageManager
	f := &PageManager{
		name:        filepath.Join(dir, name),
		fp:          fp,
		pageHeaders: make([]*pageHeader, 0),
		pids:        new(autoPageID),
	}
	// call load
	err = f.load()
	if err != nil {
		return nil, err
	}
	// return page PageManager
	return f, nil
}

// load files out the meta ([]*pageHeader) slice
// in the page pageManagerFile for easier page handling
func (f *PageManager) load() error {
	// get PageManager size info
	fi, err := f.fp.Stat()
	if err != nil {
		return err
	}
	// if this is the first run
	// not much to do, just return
	if fi.Size() < 1 {
		return nil
	}
	// otherwise, there should be
	// page headers we can load in
	for {
		// read page header data
		var h pageHeader
		_, err := readPageHeader(f.fp, &h)
		// check for an error
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return err
		}
		// add to list
		f.pageHeaders = append(f.pageHeaders, &h)
		// increment page id's
		f.pids.getNewPageID()
		// check for free pageHeaders
		if h.PageIsFree() {
			f.freePages++
		}
	}
	// seek back to the start
	_, err = f.fp.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

// getPagePosition calculates the page position based
// on the pageID provided
func getPagePosition(pid uint32) int64 {
	return int64(align(int(pid*pageSize), pageSize-1))
}

// AllocatePage allocates and returns a new page. The
// newly allocated page is not persisted unless a call
// to WritePage is made
func (f *PageManager) AllocatePage() *page {
	// generate new atomic page id
	pid := f.pids.getNewPageID()
	// create and return a new page
	return NewPage(pid)
}

// GetFreeOrAllocate attempts to find a free page (a
// page that is not in use that can be reused) and if
// one cannot be found, it will allocate and return a
// new one. Any alterations to the returned page are
// not persisted unless a call to WritePage is made
func (f *PageManager) GetFreeOrAllocate() *page {
	// first check the free page count
	if f.freePages > 0 {
		// looks like we indeed have some free pageHeaders, so
		// let's iterate all the page headers that the
		// PageManager has and try to find a free one
		for _, h := range f.pageHeaders {
			// checking if page is free
			if h.PageIsFree() {
				// found one, return it!
				p, err := f.ReadPage(h.pageID)
				if err != nil {
					// something went wrong
					panic("get free or allocate: " + err.Error())
				}
				// we should be in the clear to decrement
				// the freePages counter, and return our
				// found page
				f.freePages--
				return p
			}
		}
	}
	// otherwise, found no free pages in out freePages
	// count, so we must create and return a fresh one,
	// but first we need a fresh pageID
	pid := f.pids.getNewPageID()
	// create and return a new page with our fresh pageID
	return NewPage(pid)
}

// ReadPage attempts to read the page located at the
// offset calculated by the provided pageID. It returns
// an error if a page could not be located
func (f *PageManager) ReadPage(pid uint32) (*page, error) {
	// calc page offset in PageManager
	offset := getPagePosition(pid)
	// read data into new page
	p, err := readPageAt(f.fp, offset)
	if err != nil {
		// page not found
		return nil, ErrPageNotFound
	}
	// otherwise, return new page
	return p, nil
}

// ReadPages attempts to read the pages located at the
// offset calculated by the provided pageID. It returns
// an error if a page could not be located
func (f *PageManager) ReadPages(pid uint32) ([]*page, error) {
	// calc page offset in PageManager
	offset := getPagePosition(pid)
	// read data into new page
	p, err := readPageAt(f.fp, offset)
	if err != nil {
		// page not found
		return nil, ErrPageNotFound
	}
	// check to ensure it's an overflow page
	if p.header.hasOverflow == 0 {
		// not an overflow page
		return nil, ErrPageIsNotOverflow
	}
	// otherwise, initialize a new set of pages
	// to append the read overflow pages into
	var pages []*page
	pages = append(pages, p)
	for p.header.nextPageID > 0 {
		// calc page offset in PageManager
		offset = getPagePosition(p.header.nextPageID)
		// read data into new page
		p, err = readPageAt(f.fp, offset)
		if err != nil {
			// page not found
			return nil, ErrPageNotFound
		}
		// append it to the page set
		pages = append(pages, p)
	}
	// finally, return page set
	return pages, nil
}

// WritePage writes the provided page to the underlying PageManager
// on disk. If something goes wrong it returns a non-nil error
func (f *PageManager) WritePage(p *page) error {
	// calc page offset in PageManager
	offset := getPagePosition(p.header.pageID)
	// write provided page to PageManager
	_, err := writePageAt(f.fp, p, offset)
	if err != nil {
		// something happened
		return ErrWritingPage
	}
	// otherwise, we're good
	return nil
}

// WritePages writes the provided pages to the underlying PageManager
// on disk. If something goes wrong it returns a non-nil error
func (f *PageManager) WritePages(ps []*page) error {
	// iterate the pages
	for i := range ps {
		// page at index i
		p := ps[i]
		// calc page offset in PageManager
		offset := getPagePosition(p.header.pageID)
		// write provided page to PageManager
		_, err := writePageAt(f.fp, p, offset)
		if err != nil {
			// something happened
			return ErrWritingPage
		}
		// otherwise, we're good
	}
	return nil
}

// DeletePage marks the page with the matching pageID provided
// as "free" and writes zeros to the underlying page on disk
func (f *PageManager) DeletePage(pid uint32) error {
	// calc page offset in PageManager
	offset := getPagePosition(pid)
	// write zeros to the page found
	// at "offset" on the underlying
	// storage PageManager
	_, err := deletePageAt(f.fp, pid, offset)
	if err != nil {
		// something happened
		return ErrDeletingPage
	}
	// update the page in the
	// slotted PageManager's page cache
	for i := range f.pageHeaders {
		if f.pageHeaders[i].pageID == pid {
			// reset this matching page header
			// in the page cache to default values
			f.pageHeaders[i].freeSpaceLower = pageHeaderSize
			f.pageHeaders[i].freeSpaceUpper = pageSize
			f.pageHeaders[i].slotCount = 0
			f.pageHeaders[i].freeSlotCount = 0
		}
	}
	// otherwise, we're good
	return nil
}

// GetFreePageIDs returns a list of any
// page id's that are marked "free"
func (f *PageManager) GetFreePageIDs() []uint32 {
	// create new empty set of page id's
	var pids []uint32
	// range the page headers in
	// the slotted files page cache
	for i := range f.pageHeaders {
		if f.pageHeaders[i].PageIsFree() {
			pids = append(pids, f.pageHeaders[i].pageID)
		}
	}
	// return any free page id's found
	return pids
}

func (f *PageManager) Range(start uint32, fn func(rid *RecordID) bool) {
	for i := range f.pageHeaders {

		if f.pageHeaders[i].hasOverflow == 0 {
			break
		}
		p, err := f.ReadPage(f.pageHeaders[i].pageID)
		if err != nil {
			panic("something happend")
		}
		p.Range(fn)
		p, err = f.ReadPage(p.header.nextPageID)
		if err != nil {
			panic("something happend while reading next page")
		}
	}
}

// PageCount returns the total number of pages
// in the PageManager (including "free" pages)
func (f *PageManager) PageCount() int {
	return len(f.pageHeaders)
}

// Close closes the underlying
// PageManager, after flushing any
// buffers to disk.
func (f *PageManager) Close() error {
	err := f.fp.Close()
	if err != nil {
		return err
	}
	return nil
}

// Remove is somewhat of a helper
// function to facilitate easier
// removing of a PageManager
func Remove(path string) error {
	// sanitize path
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	// sanitize slashes
	path = filepath.ToSlash(path)
	// remove PageManager
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func (f *PageManager) grow(sizeToGrow int64) {
	fi, err := os.Stat(f.name)
	if err != nil {
		panic(err)
	}
	sizeToGrow += fi.Size()
	err = f.fp.Truncate(sizeToGrow)
	if err != nil {
		panic(err)
	}
}
