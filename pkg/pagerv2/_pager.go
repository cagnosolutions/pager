package pagerv2

import (
	"errors"
	"fmt"
	"os"
)

const pageSize = 4 << 10

type page []byte

type PageCache struct {
	entries map[uint32]uint32
	free    bitset
	data    []byte
	pages   int
	file    *os.File
	used    int
}

func NewPageCache(file *os.File, pages int) *PageCache {
	return &PageCache{
		entries: make(map[uint32]uint32),
		free:    bitset(0),
		data:    make([]byte, pages*pageSize),
		pages:   pages,
		file:    file,
	}
}

// 1: Attempt to access a page.
// 2: If the page is valid (in memory) then continue processing as normal.
// 3: If the page is invalid (not in memory) then we may need to attempt to
//    page in the page from the disk. If this is an on-demand style paging
//    system then you would issue a page-fault trap (panic.)
// 4: Next, we must check if the memory reference is a valid reference to an
//    actual location on secondary memory (disk.) If not then issue an illegal
// 	  memory access error. otherwise, go to step 5.
// 5: Read in (page in) the required page from disk. If this is an on-demand
//    style paging system then this is also the point at which you would
//    restart or continue the instruction that was interrupted by the page
//    -fault trap.

func (pc *PageCache) getFree() int {
	free := -1
	for i := uint(0); i < 63; i++ {
		if pc.free.get(i) != 0 {
			free = int(i)
			pc.free.unset(i)
			break
		}
	}
	return free
}

func (pc *PageCache) SwapIn(pid uint32) {
	// if this page is not in memory, swap it in from the disk to the cache
}

func (pc *PageCache) SwapOut(pid uint32) {
	// if this page is in memory, swap it out from the cache to the disk
}

func (pc *PageCache) NewPage() page {
	off := pc.used
	return pc.data[off : off+pageSize]
}

func (pc *PageCache) ReadPage(pid uint32) (page, error) {
	// check to see if the page is already loaded into the cache.
	off, found := pc.entries[pid]
	if !found {
		// could not be found in memory, try to acquire a new one.
		return nil, errors.New("no data, cache miss")
	}
	return pc.data[off : off+pageSize], nil
}

func (pc *PageCache) WritePage(p page) (uint32, error) {
	free := pc.getFree()
	if free == -1 {
		return 0, errors.New("no free pages left!")
	}
	off := free * pageSize
	copy(pc.data[off:off+pageSize], p)
	return uint32(off), nil
}

func (pc *PageCache) DelPage(pid uint32) error {
	off, found := pc.entries[pid]
	if !found {
		return errors.New("no data, cache miss")
	}
	delete(pc.entries, pid)
	pc.free.set(uint(off))
	return nil
}

type bitset uint64

func (bs *bitset) set(i uint) {
	*bs |= 1 << (i & (63))
}

func (bs *bitset) unset(i uint) {
	*bs &^= 1 << (i & (63))
}

func (bs *bitset) get(i uint) uint {
	return uint(*bs) & (1 << (i & (63)))
}

func (bs bitset) print() {
	fmt.Printf("%.64b\n", bs)
}
