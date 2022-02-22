package pager

import "encoding/binary"

var bindata = binary.LittleEndian

func setFreshHeader(p []byte, pid uint32) {
	if len(p) != pageHeaderSize {
		panic(ErrBadPageSize)
	}
	bindata.PutUint32(p[0:4], pid)              // pageID
	bindata.PutUint32(p[4:8], 0)                // nextPageID
	bindata.PutUint32(p[8:12], 0)               // prevPageID
	bindata.PutUint16(p[12:14], pageHeaderSize) // freeSpaceLower
	bindata.PutUint16(p[14:16], pageSize)       // freeSpaceUpper
	bindata.PutUint16(p[16:18], 0)              // slotCount
	bindata.PutUint16(p[18:20], 0)              // freeSlotCount
	bindata.PutUint16(p[20:22], 0)              // hasOverflow
	bindata.PutUint16(p[22:24], 0)              // reserved
}

/*
	Setter functions for getting values from raw pages
*/
func setPageID(b []byte, pageID uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	bindata.PutUint32(b, pageID)
}

func setNextPageID(b []byte, nextPageID uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	bindata.PutUint32(b, nextPageID)
}

func setPrevPageID(b []byte, prevPageID uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	bindata.PutUint32(b, prevPageID)
}

func setFreeSpaceLower(b []byte, freeSpaceLower uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, freeSpaceLower)
}

func setFreeSpaceUpper(b []byte, freeSpaceUpper uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, freeSpaceUpper)
}

func setSlotCount(b []byte, slotCount uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, slotCount)
}

func setFreeSlotCount(b []byte, freeSlotCount uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, freeSlotCount)
}

func setHasOverflow(b []byte, hasOverflow uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, hasOverflow)
}

func setReserved(b []byte, reserved uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	bindata.PutUint16(b, reserved)
}

/*
	Getter functions for getting values from raw pages
*/
func getPageID(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint32(b)
}

func getNextPageID(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint32(b)
}

func getPrevPageID(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint32(b)
}

func getFreeSpaceLower(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

func getFreeSpaceUpper(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

func getSlotCount(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

func getFreeSlotCount(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

func getHasOverflow(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

func getReserved(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b)
}

// slot layout below
// itemID     uint16
// itemStatus uint16
// itemOffset uint16
// itemLength uint16

func getSlotNOffset(slotNumber int) int {
	return pageHeaderSize + slotNumber*pageSlotSize
}

func getSlotID(b []byte, slotNumber int) uint16 {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryID
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b[n : n+2])
}

func setSlotID(b []byte, slotNumber int, slotID uint16) {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryID
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	bindata.PutUint16(b[n:n+2], slotID)
}

func getSlotStatus(b []byte, slotNumber int) uint16 {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryStatus
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b[n : n+2])
}

func setSlotStatus(b []byte, slotNumber int, slotStatus uint16) {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryStatus
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	bindata.PutUint16(b[n:n+2], slotStatus)
}

func getSlotOffset(b []byte, slotNumber int) uint16 {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryOffset
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b[n : n+2])
}

func setSlotOffset(b []byte, slotNumber int, slotOffset uint16) {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryOffset
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	bindata.PutUint16(b[n:n+2], slotOffset)
}

func getSlotLength(b []byte, slotNumber int) uint16 {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryLength
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	return bindata.Uint16(b[n : n+2])
}

func setSlotLength(b []byte, slotNumber int, slotLength uint16) {
	n := getSlotNOffset(slotNumber)
	n += offSlotEntryLength
	_ = b[n] // bounds check hint to compiler; see golang.org/issue/14808
	bindata.PutUint16(b[n:n+2], slotLength)
}
