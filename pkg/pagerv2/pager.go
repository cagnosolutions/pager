package pagerv2

import (
	"io"
	"os"
)

const (
	pageSize = 4 << 10
)

type Pager struct {
	file          *os.File
	cache         *lru
	data          []byte
	usedNumPages  int
	dirtyNumPages int
	maxNumPages   int
}

func NewPager(path string, pages int) *Pager {
	fp, err := OpenFile(path)
	if err != nil {
		panic(err)
	}
	p := &Pager{
		file:          fp,
		cache:         newLRU(pages),
		data:          make([]byte, pageSize*pages),
		usedNumPages:  0,
		dirtyNumPages: 0,
		maxNumPages:   pages,
	}
	err = p.load()
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Pager) load() error {
	// first we get the file size information
	fi, err := p.file.Stat()
	if err != nil {
		return err
	}
	// if this is the first run and the file is empty, then we should
	// initialize the cache for the first p.pages entries and return
	if fi.Size() < 1 {
		// add p.pages entries to the cache to utilize
		for i := 0; i < p.maxNumPages; i++ {
			p.cache.set(uint32(i), int64(i*pageSize))
		}
		// and then return
		return nil
	}
	// otherwise, there should be some pages we can load in (we are simply
	// looking for a page offset count at this point)
	for i := uint32(0); i < uint32(p.maxNumPages) && err != io.EOF; i++ {
		// calculate current page offset
		off := int64(i * pageSize)
		// read page into pager
		_, err := p.file.ReadAt(p.data[off:off+pageSize], off)
		// check for an error
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return err
		}
		// add to cache
		p.cache.set(i, off)
		// update pager info
		p.usedNumPages++
	}
	// done
	return nil
}

// INDEXING STRUCTURE: https://go.dev/play/p/8lTKeR4fLYj

func (p *Pager) evict(numPages int) {
	// evict numPages and add to free set
}

func (p *Pager) free() (uint32, int64) {
	// return a free page
	return 0, 0
}

func (p *Pager) getPage(off int64) ([]byte, error) {
	// error check offset
	if off > int64(len(p.data)) || off < 0 {
		// encountered error, return nil and illegal access
		return nil, ErrIllegalPageAccess
	}
	// page align offset (in case its off)
	off = off&^pageSize - 1
	// return page, and nil error
	return p.data[off : off+pageSize], nil
}

func (p *Pager) Read(pid uint32) ([]byte, error) {
	// look for the page in memory
	off, found := p.cache.get(pid)
	if found {
		// if we find it then return it (this counts as a cache hit)
		return p.getPage(off)
	}
	// since we did not find the page cache in memory, we should check to
	// make sure our page buffer is not full
	if p.usedNumPages == p.maxNumPages {
		// page cache is full, check the dirty page count
		if p.dirtyNumPages == p.maxNumPages {
			// we are full and all the pages are dirty which means we need
			// to flush our pages
			err := p.Sync()
			if err != nil {
				return nil, err
			}
		}
		// we may have had to sync, but either way we have to evict
		p.evict(4)
	}
	// acquire a free page
	pid, off = p.free()
	// and since we did not find it in the cache, we need to attempt to
	// read it off the disk and cache it (this counts as a cache miss)
	_, err := p.file.ReadAt(p.data[off:off+pageSize], off)
	if err != nil {
		return nil, err
	}
	// next, cache the newly read page
	p.cache.set(pid, off)
	// and finally, return the page data
	return p.getPage(off)
}

func (p *Pager) Write(d []byte, pid uint32) error {
	// 1) look for page in memory
	// 2) if we find it proceed, and if not we need to figure out what to do
	// 3) check to make sure the page has room for the record, and if not
	//    then return page is full or record too large error
	// 4) copy the data to the page
	return nil
}

func (p *Pager) GetFreePageID() uint32 {
	// TODO implement me
	panic("implement me")
}

func (p *Pager) GetFreeRecordID(pid uint32, rsize int) (uint16, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Pager) ReadRecord(pid uint32, rid uint16) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (p *Pager) WriteRecord(r []byte, pid uint32, rid uint16) error {
	// TODO implement me
	panic("implement me")
}

func (p *Pager) Sync() error {
	// TODO implement me
	panic("implement me")
}
