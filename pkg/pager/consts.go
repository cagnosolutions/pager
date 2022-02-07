package pager

const (
	pageSize       = 8 << 10 // 8 KB
	pageHeaderSize = 24      // 24 bytes
	pageSlotSize   = 8       // 8 bytes
	MinRecordSize  = pageSlotSize
	MaxRecordSize  = pageSize - pageHeaderSize - pageSlotSize
)

const (
	itemStatusFree uint16 = iota
	itemStatusUsed
)

func align(n int, size int) int {
	return (n + size) &^ size
}
