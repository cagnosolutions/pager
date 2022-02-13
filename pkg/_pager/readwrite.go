package _pager

import (
	"encoding/binary"
	"io"
)

/*
	pageID         uint32
	nextPageID     uint32
	prevPageID     uint32
	freeSpaceLower uint16
	freeSpaceUpper uint16
	slotCount      uint16
	freeSlotCount  uint16
	hasOverflow    uint16
	reserved       uint16
*/

// PageReader reads a pageSize of bytes into p starting at
// offset off in the underlying input source. It returns
// the number of bytes read and any error encountered
//
// If ReadPage is reading from an input source with a seek
// offset, ReadPage should not affect nor be affected by
// the underlying seek offset.
type PageReader interface {
	ReadPage(p []byte, off int64) (n int, err error)
}

// PageWriter writes a pageSize of bytes from p to the under-
// lying data stream at offset off. It returns the number of
// bytes written from p and any error encountered that caused
// the write to stop early. WritePage must return a non-nil
// error if it returns n < a pageSize.
//
// If WritePage is writing to a destination with a seek offset,
// WritePage should not affect nor be affected by the underlying
// seek offset.
type PageWriter interface {
	WritePage(p []byte, off int64) (n int, err error)
}

func readPageHeader(r io.ReadSeeker, h *pageHeader) (int, error) {
	// make header buffer to read data into
	buf := make([]byte, pageHeaderSize)
	// read the header from the underlying reader into the buffer
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	// decode PageID
	h.pageID = binary.LittleEndian.Uint32(buf[0:4])
	// decode nextPageID
	h.nextPageID = binary.LittleEndian.Uint32(buf[4:8])
	// decode prevPageID
	h.prevPageID = binary.LittleEndian.Uint32(buf[8:12])
	// decode freeSpaceLower
	h.freeSpaceLower = binary.LittleEndian.Uint16(buf[12:14])
	// decode freeSpaceUpper
	h.freeSpaceUpper = binary.LittleEndian.Uint16(buf[14:16])
	// decode slotCount
	h.slotCount = binary.LittleEndian.Uint16(buf[16:18])
	// decode freeSlotCount
	h.freeSlotCount = binary.LittleEndian.Uint16(buf[18:20])
	// decode hasOverflow
	h.hasOverflow = binary.LittleEndian.Uint16(buf[20:22])
	// decode reserved
	h.reserved = binary.LittleEndian.Uint16(buf[22:24])
	// seek to the start of the next Page header
	nn, err := r.Seek(int64(pageSize-n), io.SeekCurrent)
	if err != nil {
		return int(nn) + n, err
	}
	// return bytes read, and error
	return int(nn) + n, nil
}

func decodePageHeader(b []byte, h *pageHeader) int {
	// get offset to decode Page header directly into Page data
	var n int
	// decode PageID
	h.pageID = binary.LittleEndian.Uint32(b[n : n+4])
	n += 4
	// decode nextPageID
	h.nextPageID = binary.LittleEndian.Uint32(b[n : n+4])
	n += 4
	// decode prevPageID
	h.prevPageID = binary.LittleEndian.Uint32(b[n : n+4])
	n += 4
	// decode freeSpaceLower
	h.freeSpaceLower = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// decode freeSpaceUpper
	h.freeSpaceUpper = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// decode slotCount
	h.slotCount = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// decode freeSlotCount
	h.freeSlotCount = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// decode hasOverflow
	h.hasOverflow = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// decode reserved
	h.hasOverflow = binary.LittleEndian.Uint16(b[n : n+2])
	n += 2
	// return
	return n
}

func encodePageHeader(b []byte, h *pageHeader) int {
	// get offset to encode Page header directly into Page data
	var n int
	// encode PageID
	binary.LittleEndian.PutUint32(b[n:n+4], h.pageID)
	n += 4
	// encode nextPageID
	binary.LittleEndian.PutUint32(b[n:n+4], h.nextPageID)
	n += 4
	// encode prevPageID
	binary.LittleEndian.PutUint32(b[n:n+4], h.prevPageID)
	n += 4
	// encode freeSpaceLower
	binary.LittleEndian.PutUint16(b[n:n+2], h.freeSpaceLower)
	n += 2
	// encode freeSpaceUpper
	binary.LittleEndian.PutUint16(b[n:n+2], h.freeSpaceUpper)
	n += 2
	// encode slotCount
	binary.LittleEndian.PutUint16(b[n:n+2], h.slotCount)
	n += 2
	// encode freeSlotCount
	binary.LittleEndian.PutUint16(b[n:n+2], h.freeSlotCount)
	n += 2
	// encode hasOverflow
	binary.LittleEndian.PutUint16(b[n:n+2], h.hasOverflow)
	n += 2
	// encode reserved
	binary.LittleEndian.PutUint16(b[n:n+2], h.reserved)
	n += 2
	// return bytes encoded
	return n
}

func readPage(r io.Reader, p *Page) (int, error) {
	// init Page data
	p.data = make([]byte, pageSize)
	// read Page data
	nn, err := r.Read(p.data)
	if err != nil {
		return -1, err
	}
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
	// return bytes read
	return nn, nil
}

func readPageAt(r io.ReaderAt, offset int64) (*Page, error) {
	// init new Page
	p := new(Page)
	// init new Page data
	p.data = make([]byte, pageSize)
	// read Page data into Page from the
	// underlying pageManagerFile at the offset provided
	_, err := r.ReadAt(p.data, offset)
	if err != nil {
		return nil, err
	}
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
	// return read Page
	return p, nil
}

func writePage(w io.Writer, p *Page) (int, error) {
	// encode Page header
	n := encodePageHeader(p.data[0:pageHeaderSize], p.header)
	// encode Page slots
	for i := range p.slots {
		// encode slot item prefix
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemID)
		n += 2
		// encode slot item status
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemStatus)
		n += 2
		// encode slot item offset
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemOffset)
		n += 2
		// encode slot item length
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemLength)
		n += 2
	}
	// write Page data
	nn, err := w.Write(p.data)
	if err != nil {
		return nn, err
	}
	// return bytes written and possible error
	return nn, nil
}

func writePageAt(w io.WriterAt, p *Page, offset int64) (int, error) {
	// encode Page header
	n := encodePageHeader(p.data[0:pageHeaderSize], p.header)
	// encode Page slots
	for i := range p.slots {
		// encode slot item prefix
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemID)
		n += 2
		// encode slot item status
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemStatus)
		n += 2
		// encode slot item offset
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemOffset)
		n += 2
		// encode slot item length
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemLength)
		n += 2
	}
	// write Page data to the underlying
	// pageManagerFile at the offset provided
	nn, err := w.WriteAt(p.data, offset)
	if err != nil {
		return nn, err
	}
	// return bytes written and possible error
	return nn, nil
}

func deletePageAt(w io.WriterAt, pid uint32, offset int64) (int, error) {
	// create a new "empty" Page
	p := NewPage(pid)
	// encode Page header
	n := encodePageHeader(p.data[0:pageHeaderSize], p.header)
	// encode Page slots
	for i := range p.slots {
		// encode slot item prefix
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemID)
		n += 2
		// encode slot item status
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemStatus)
		n += 2
		// encode slot item offset
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemOffset)
		n += 2
		// encode slot item length
		binary.LittleEndian.PutUint16(p.data[n:n+2], p.slots[i].itemLength)
		n += 2
	}
	// write Page data to the underlying
	// pageManagerFile at the offset provided
	nn, err := w.WriteAt(p.data, offset)
	if err != nil {
		return nn, err
	}
	// return bytes written and possible error
	return nn, nil
}
