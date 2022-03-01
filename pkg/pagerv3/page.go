package pagerv3

const pageSize = 16

type Page struct {
	id       PageID
	pinCount int
	isDirty  bool
	data     [pageSize]byte
}
