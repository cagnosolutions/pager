package pager

const (
	// used in Page
	pageSize       = 8 << 10 // 8 KB
	pageHeaderSize = 24      // 24 bytes
	pageSlotSize   = 8       // 8 bytes
	MinRecordSize  = pageSlotSize
	MaxRecordSize  = pageSize - pageHeaderSize - pageSlotSize

	// used in PageBuffer
	defaultBufferedPageCount = 8
)

const (
	itemStatusFree uint16 = iota
	itemStatusUsed
)

func align(n int, size int) int {
	return (n + size) &^ size
}
