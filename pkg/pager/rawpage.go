package pager

import (
	"bytes"
	"encoding/binary"
	"sort"
	"sync"
)



// page is a raw page sized []byte
type page []byte

// initPage initializes a fresh page in the page space provided and encodes the
// header with the page id provided
func initPage(p []byte, pid uint32) {
	// encode page header
	setFreshHeader(p[:pageHeaderSize], pid)
	// and we're done!
}

// freeSpace returns the free space left in the page
func (p page) freeSpace() uint16 {
	return getFreeSpaceUpper(p) - getFreeSpaceLower(p)
}

// getAvailableSlot returns the offset of a free page slot if there is one
// already existing that can be used. If a free page slot cannot be found it
// adds a new slot to the current page and returns the slot number of the slot
// (in the page) of the found or inserted page slot that can be used for
// inserting a new record.
func (p page) getAvailableSlot(recordSize uint16) int {
	// first we will check the page header to see if the freeSlotCount is
	// reporting any free slots
	if p.getFreeSlotCount() > 0 {
		// it looks like we might have one so let's begin scanning the existing
		// set of slots and check for any that are currently marked as free.
		// as a note, `i` will be our current index within the page and c will
		// act as out slot counter. we will start `i` out at the end of the
		// header and increment it pageSlotSize (8 bytes) with every iteration.
		for slot := 0; slot < int(p.getSlotCount()); slot++ {
			// now we must read the status of the current slot
			if getSlotStatus(p, slot) == itemStatusFree {
				// looks like we found a free slot. we can now call our use
				// free page slot record to recycle a page slot. it takes the
				// offset of the head of the slot to use, along with the record
				// size. it will recycle the slot and return the slot offset.
				return p.useFreePageSlotRecord(slot, recordSize)
			}
		}
		// we should NEVER get here
		panic("if you see this, look in rawpage.go in getAvailableSlot()")
	}
	// otherwise, looks like we aren't reporting that we have any existing page
	// slots free, so we should just add a new page slot record and be done.
	return p.addNewPageSlotRecord(recordSize)
}

// useFreePageSlotRecord uses an existing page slot record provided. it attempts
// to use the same record offset (if it will fit) otherwise, it will find
// another location in the page and update the header accordingly.
func (p page) useFreePageSlotRecord(slotNum int, recordSize uint16) int {
	// no need to increment the slot counter, but we do need to decrement the
	// free slot count because we are now using one of those free slots
	freeSlotCount := p.getFreeSlotCount()
	freeSlotCount--
	p.setFreeSlotCount(freeSlotCount)
	// let's check to see if the last record offset had enough space to fit
	// the current record in
	if recordSize <= getSlotLength(p, slotNum) {
		// it will fit, nice! let us update the slot to fit the new record
		// information--remember the only things that will change are the
		// slot entry status, and the slot entry length.
		setSlotStatus(p, slotNum, itemStatusUsed)
		setSlotLength(p, slotNum, recordSize)
		// we are not finished, just return the slot number
		return slotNum
	}
	// now, if we are here it means that the last record was not large enough
	// to house this record data, so first we just need to allocate a whole new
	// offset to store the record that is the size of the new record. to do
	// that, we will adjust the free space upper bound
	freeSpaceUpper := p.getFreeSpaceUpper()
	freeSpaceUpper -= recordSize
	p.setFreeSpaceUpper(freeSpaceUpper)
	// we do not need to raise the free space lower bound because we are not
	// adding a new slot--but now we need to update the slot entry status,
	// entry offset and entry length...
	setSlotStatus(p, slotNum, itemStatusUsed)
	setSlotOffset(p, slotNum, p.getFreeSpaceUpper())
	setSlotLength(p, slotNum, recordSize)
	// we should be all done, return the slot number
	return slotNum
}

// addNewPageSlotRecord appends a new page slot to the slots list within the
// page, and it updates the page header, incrementing the slot count, growing
// the free space lower bound and shrinking the free space upper bound. it
// returns the page slot number of the newly added page slot.
func (p page) addNewPageSlotRecord(recordSize uint16) int {
	// first we increment the slot count
	slotCount := p.getSlotCount()
	slotCount++
	p.setSlotCount(slotCount)
	// next, we raise the free space lower boundary because we are now adding a
	// new slot
	freeSpaceLower := p.getFreeSpaceLower()
	freeSpaceLower += pageSlotSize
	p.setFreeSpaceLower(freeSpaceLower)
	// then, we must lower the free space upper bound because we are adding
	// the record data
	freeSpaceUpper := p.getFreeSpaceUpper()
	freeSpaceUpper -= recordSize
	p.setFreeSpaceUpper(freeSpaceUpper)
	// finally, before we return we need to create a new page slot entry
	// recording the byte offset where the record will be copied to within
	// the page along with the length of the record.
	// but first, we calculate the slot number.
	slotNum := int(p.getSlotCount() - 1)
	// then we write the new slot to the page
	p.addSlot(slotNum, recordSize)
	// return the last page slot we entered
	return slotNum
}

func (p page) addSlot(slotNum int, recordSize uint16) {
	setSlotID(p, slotNum, uint16(slotNum))
	setSlotStatus(p, slotNum, itemStatusUsed)
	setSlotOffset(p, slotNum, p.getFreeSpaceUpper())
	setSlotLength(p, slotNum, recordSize)
}

// itemBounds returns the beginning and ending offset positions for the
// location of this item within the page
func (p page) slotEntryBounds(slotNum int) (uint16, uint16) {
	offset := getSlotOffset(p, slotNum)
	length := getSlotLength(p, slotNum)
	return offset, offset + length
}

// writeRecord attempts to write a new record to a raw page and return an id.
// If something fails it will return an error and an empty id
func (p page) writeRecord(rec []byte) (uint32, error) {
	// check record to make sure it is will fit
	recordSize := len(rec)
	if recordSize < MinRecordSize {
		return 0, ErrMinRecordSize
	}
	if recordSize > MaxRecordSize {
		return 0, ErrMaxRecordSize
	}
	if recordSize > int(p.freeSpace()) {
		return 0, ErrNoMoreRoomInPage
	}
	// if we are here, this means the record is fine now we must look for free
	// slot to write into. get a fresh (or used free one, if there are any) and
	// update the page header accordingly
	sid := p.getAvailableSlot(uint16(recordSize))
	// get the new record offsets
	beg, end := p.slotEntryBounds(sid)
	// copy the record to the Page
	copy(p[beg:end], rec)
	// before we return (this does not affect the page slot id) we should sort
	// the slot pointers, so all the record pointers are in the proper order.
	p.sortSlotsByRecordPrefix()
	// all went well, return the newly added
	// records' id along with a nil error
	return &RecordID{
		PageID: p.header.pageID,
		SlotID: s.itemID,
	}, nil
}

// Len is here to satisfy the sort interface for
// sorting the Page slots by the record prefix
func (p page) Len() int {
	return int(p.getSlotCount())
}

// Swap is here to satisfy the sort interface for
// sorting the Page slots by the record prefix
func (p page) Swap(i, j int) {
	// get slot `i` and slot `j`
	slotI := p.getWholePageSlot(i)
	slotJ := p.getWholePageSlot(j)
	// encode slot `i` and slot `j` into uint64's
	slotIu64 := binary.LittleEndian.Uint64(slotI)
	slotJu64 := binary.LittleEndian.Uint64(slotJ)
	// swap slots `i` and `j`
	slotIu64, slotJu64 = slotJu64, slotIu64
	// copy slot `i` into our slot pool
	copy(sp, slotI)
	// copy slot `j` into slot `i`
	copy(slotI, slotJ)
	// copy slot `i` (slot pool) into slot `j`
	copy(slotJ, sp)
	// p.slots[i], p.slots[j] = p.slots[j], p.slots[i]
}

// Less is here to satisfy the sort interface for
// sorting the Page slots by the record prefix
func (p page) Less(i, j int) bool {
	p.slotEntryBounds(i)
	ipre, _ := p.slotEntryBounds(i)
	jpre, _ := p.slotEntryBounds(j)
	return bytes.Compare(p[ipre:ipre+8], p[jpre:jpre+8]) < 0
}

// https://go.dev/play/p/Tx0PwDL1UW1


// sortSlotsByRecordPrefix is a wrapper for sorting the page slots by the
// record prefix
func (p page) sortSlotsByRecordPrefix() {
	ss := p.slotsToSlotSet()
	sort.Stable(slots)
	p.slotSetToSlots(ss)
}

func (p page) slotsToSlotSet() slotSet {
	set := make(slotSet, p.getSlotCount())
	for i := range set {
		set[i].slotID
		set[i].

	}
}

func (p page) getPageID() uint32 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offPageID
	return getPageID(p[n : n+4])
}

func (p page) getNextPageID() uint32 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offNextPageID
	return getNextPageID(p[n : n+4])
}

func (p page) setNextPageID(nextPageID uint32) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offNextPageID
	setNextPageID(p[n:n+4], nextPageID)
}

func (p page) getPrevPageID() uint32 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offPrevPageID
	return getPrevPageID(p[n : n+4])
}

func (p page) setPrevPageID(prevPageID uint32) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offPrevPageID
	setPrevPageID(p[n:n+4], prevPageID)
}

func (p page) getFreeSpaceLower() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offFreeSpaceLower
	return getFreeSpaceLower(p[n : n+2])
}

func (p page) setFreeSpaceLower(freeSpaceLower uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offFreeSpaceLower
	setFreeSpaceLower(p[n:n+2], freeSpaceLower)
}

func (p page) getFreeSpaceUpper() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offFreeSpaceUpper
	return getFreeSpaceUpper(p[n : n+2])
}

func (p page) setFreeSpaceUpper(freeSpaceUpper uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offFreeSpaceUpper
	setFreeSpaceUpper(p[n:n+2], freeSpaceUpper)
}

func (p page) getSlotCount() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offSlotCount
	return getSlotCount(p[n : n+2])
}

func (p page) setSlotCount(slotCount uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offSlotCount
	setSlotCount(p[n:n+2], slotCount)
}

func (p page) getFreeSlotCount() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offFreeSlotCount
	return getFreeSlotCount(p[n : n+2])
}

func (p page) setFreeSlotCount(freeSlotCount uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offFreeSlotCount
	setFreeSlotCount(p[n:n+2], freeSlotCount)
}

func (p page) getHasOverflow() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offHasOverflow
	return getHasOverflow(p[n : n+2])
}

func (p page) setHasOverflow(hasOverflow uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offHasOverflow
	setHasOverflow(p[n:n+2], hasOverflow)
}

func (p page) getReserved() uint16 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	n := offReserved
	return getReserved(p[n : n+2])
}

func (p page) setReserved(reserved uint16) {
	// early bounds check to guarantee safety of writes below
	_ = p[pageSize-1]
	n := offReserved
	setReserved(p[n:n+2], reserved)
}

func (p page) getPageHeader() []byte {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	return p[:pageHeaderSize]
}

func (p page) getWholePageSlot(slotNum int) []byte {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	slotOffset := getSlotNOffset(slotNum)
	return p[slotOffset : slotOffset+pageSlotSize]
}

func (p page) slotToU64(slotNum int) uint64 {
	// bounds check hint to compiler; see golang.org/issue/14808
	_ = p[pageSize-1]
	slotOffset := getSlotNOffset(slotNum)
	return bindata.Uint64(p[slotOffset:slotOffset+pageSlotSize])
}

func (p page) AddRecord(rec []byte) (uint, error) {
	return 0, nil
}

func (p page) GetRecord(rec []byte) (uint, error) {
	return 0, nil
}

func (p page) DelRecord(rec []byte) (uint, error) {
	return 0, nil
}

func (p page) RangeRecords(fn func(rec []byte) bool) {

}

func (p page) Reset() {

}

func (p page) String() string {
	return ""
}

func (p page) writeNewPageHeader(pid uint32) {
	// get offset to encode Page header directly into Page data
	var n int
	// encode PageID
	binary.LittleEndian.PutUint32(p[n:n+4], pid)
	n += 4
	// encode nextPageID
	binary.LittleEndian.PutUint32(p[n:n+4], 0)
	n += 4
	// encode prevPageID
	binary.LittleEndian.PutUint32(p[n:n+4], 0)
	n += 4
	// encode freeSpaceLower
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// encode freeSpaceUpper
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// encode slotCount
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// encode freeSlotCount
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// encode hasOverflow
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// encode reserved
	binary.LittleEndian.PutUint16(p[n:n+2], 0)
	n += 2
	// return bytes encoded
	return n
}
