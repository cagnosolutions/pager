package _pager

//	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//	NOTE: this file will most likely be removed soon...
//	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

type chunk struct {
	pages []*Page
}

func genID(cid uint32, i int) uint32 {
	return cid + uint32(i)
}

func genNextID(cid uint32, i int) uint32 {
	if i < 7 {
		return cid + uint32(i) + 1
	}
	return 0
}

func genPrevID(cid uint32, i int) uint32 {
	if i > 0 {
		return cid + uint32(i) - 1
	}
	return 0
}

func _NewChunk(cid uint32) *chunk {
	c := &chunk{
		pages: make([]*Page, 8),
	}
	for i := range c.pages {
		c.pages[i] = &Page{
			header: &pageHeader{
				pageID:         genID(cid, i),
				nextPageID:     genNextID(cid, i),
				prevPageID:     genPrevID(cid, i),
				freeSpaceLower: pageHeaderSize,
				freeSpaceUpper: pageSize,
				slotCount:      0,
				freeSlotCount:  0,
				hasOverflow:    1,
				reserved:       0,
			},
			slots: make([]*pageSlot, 0),
			data:  make([]byte, pageSize),
		}
	}
	return c
}
