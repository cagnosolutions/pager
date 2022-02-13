package pager

import (
	"io"
	"os"
)

// Cache is a page table cache [https://en.wikipedia.org/wiki/Page_table]
type Cache struct {
	pagePtrs  map[uint32]uint32
	dataCache [cacheSize]byte
	fp        *os.File
}

// ReadPage attempts to read the Page located at the
// offset calculated by the provided pageID. It returns
// an error if a Page could not be located
func (c *Cache) ReadPage(pid uint32) error {
	// find an available spot in the cache
	apid := c.getAvailable()
	// calculate the page offset
	poff := pageOffset(pid)
	// swap a new page in
	err := c.swapIn(c.fp, poff, apid)
	if err != nil {
		// Page not found
		return ErrPageNotFound
	}
	// otherwise, return new Page
	return nil
}

// getAvailable locates the first available page frame
// or if none can be found it evicts the oldest one
// note: page replacement algorithms, see LRU-K, but
// fall back on LRU
func (c *Cache) getAvailable() uint32 {
	return 0
}

// swapIn reads a page from disk into main memory
func (c *Cache) swapIn(r io.ReaderAt, poff int64, apid uint32) error {
	// get the page
	p := c.getPage(apid)
	// read page into cache from the underlying file
	_, err := r.ReadAt(p, poff)
	if err != nil {
		return err
	}
	/*
		// init Page header
		p.header = new(pageHeader)
		// decode Page header
		n := decodePageHeader(p.data[0:pageHeaderSize], p.header)
		// init Page slots
		p.slots = make([]*pageSlot, p.header.slotCount)
		// decode Page slots
		for i := range p.slots {
			// create a new pageSlot pointer
			p.slots[i] = new(pageSlot)
			// encode slot item prefix
			p.slots[i].itemID = binary.LittleEndian.Uint16(p.data[n : n+2])
			n += 2
			// encode slot item status
			p.slots[i].itemStatus = binary.LittleEndian.Uint16(p.data[n : n+2])
			n += 2
			// encode slot item offset
			p.slots[i].itemOffset = binary.LittleEndian.Uint16(p.data[n : n+2])
			n += 2
			// encode slot item length
			p.slots[i].itemLength = binary.LittleEndian.Uint16(p.data[n : n+2])
			n += 2
		}
	*/
	return nil
}

// SwapOut writes from memory to the disk
func (c *Cache) SwapOut() {

}

type page []byte

// getPage returns a pointer to the raw page
// referenced by the provided page id. It will
// return nil if it encounters and error
func (c *Cache) getPage(pid uint32) []byte {
	return nil
}

// pageOffset converts a page id into an offset address
// which is used to locate the page position on disk
func pageOffset(pid uint32) int64 {
	return int64(align(int(pid*pageSize), pageSize-1))
}

//func (c *Cache) Load(r io.ReaderAt, offset int64) {
//	// load up the page cache
//	for i := 0; i < cachePageCount; i++ {
//		_, err := r.ReadAt(c.pdata[i:i+pageSize], offset)
//		if err != nil {
//			if err == io.EOF || err == io.ErrUnexpectedEOF {
//				break
//			}
//			panic(err)
//		}
//	}
//}

// Lookup attempts to acquire a page if it is already in the in-memory cache.
// It does not read the page from disk. It will return nil if it is not found.
func (c *Cache) Lookup(pid uint32) *Page {
	return nil
}

/*
	// SQLite3 -> https://github.com/smparkes/sqlite/blob/master/src/pager.h
 	//
	// Functions used to obtain and release page references.
	int sqlite3PagerAcquire(Pager *pPager, Pgno pgno, DbPage **ppPage, int clrFlag);
	#define sqlite3PagerGet(A,B,C) sqlite3PagerAcquire(A,B,C,0)
	DbPage *sqlite3PagerLookup(Pager *pPager, Pgno pgno);
	void sqlite3PagerRef(DbPage*);
	void sqlite3PagerUnref(DbPage*);

	// Operations on page references.
	int sqlite3PagerWrite(DbPage*);
	void sqlite3PagerDontWrite(DbPage*);
	int sqlite3PagerMovepage(Pager*,DbPage*,Pgno,int);
	int sqlite3PagerPageRefcount(DbPage*);
	void *sqlite3PagerGetData(DbPage *);
	void *sqlite3PagerGetExtra(DbPage *);
*/
