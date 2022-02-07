package pager

type pageMeta struct {
	isDirty   bool
	freeSpace uint16
}

type PageBuffer struct {
	manager *PageManager
	buffer  []*page
	metas   []pageMeta
	pinned  int
}

func NewPageBufferSize(pm *PageManager, np int) (*PageBuffer, error) {
	// create page buffer
	pb := &PageBuffer{
		manager: pm,
		buffer:  make([]*page, np),
		metas:   make([]pageMeta, np),
		pinned:  0,
	}
	// call load
	err := pb.load()
	if err != nil {
		return nil, err
	}
	// pin the record that has
	// the most free space
	var pinned int
	var free uint16
	for i := range pb.metas {
		if pb.metas[i].freeSpace > free {
			free = pb.metas[i].freeSpace
			pinned = i
		}
	}
	pb.pinned = pinned
	// return it
	return pb, nil
}

func NewPageBuffer(pm *PageManager) (*PageBuffer, error) {
	return NewPageBufferSize(pm, defaultBufferedPageCount)
}

func (pb *PageBuffer) load() error {
	// get the page count that the manager has
	pageCount := len(pb.manager.pageHeaders)
	// check to see if we need to allocate any
	// new or initial pages for the page buffer
	if pageCount == 0 {
		// no pages are currently allocated
		for i := 0; i < len(pb.buffer); i++ {
			// so, we must allocate a new page
			// and append it to our buffer
			pb.buffer[i] = pb.manager.AllocatePage()
			// and now we are done with this case
		}
		// so we simply return
		return nil
	}
	// otherwise, the manager has one or more pages
	// allocated, but we must check to see if we need
	// to allocate any other pages to fill the buffer
	if pageCount > 0 && pageCount <= 8 {
		// let's get the remaining pages allocated
		for i := 0; i < len(pb.buffer)-pageCount; i++ {
			// allocate a new page and append it to
			// our page buffer page set
			pb.buffer[i] = pb.manager.AllocatePage()
			// and we are done with this case
		}
		// so we simply return
		return nil
	}
	// otherwise, we most likely have a reference to
	// more than 8 page in the manager, so (in this
	// instance) we will simply load the first 8 into
	// the page buffer
	if pageCount > 8 {
		// let's get the first eight pages from the manager
		for i := 0; i < len(pb.buffer); i++ {
			// so, we range the first 8 from the manager
			// and load them into the page buffer by first
			// getting the page id we need to read in
			pid := pb.manager.pageHeaders[i].pageID
			// read the page using the manager
			p, err := pb.manager.ReadPage(pid)
			if err != nil {
				return err
			}
			// as long as we didn't get any errors
			// we should be right as rain as they
			// say; append our page to the buffer
			// along with any and all other page
			// metadata that we need.
			pb.buffer[i] = p
			pb.metas[i] = pageMeta{
				isDirty:   false,
				freeSpace: p.header.FreeSpace(),
			}
			// and now we should be done with this
			// case, and there should be no other
			// cases, so we should...
		}
		// simply return
		return nil
	}
	return nil
}

func (pb *PageBuffer) AddRecord(r []byte) (*RecordID, error) {
	// get the record size
	size := len(r)
	// first, we do a simple error check
	if size > MaxRecordSize*len(pb.buffer) {
		return nil, ErrWritingPage
	}
	// then check to see if the record is small
	// enough to simply fit inside on page
	if size < MaxRecordSize {
		// if the record size is small enough to fit
		// in one page then write get the pinned page
		p := pb.buffer[pb.pinned]
		// and write the record to the pinned page
		rid, err := p.AddRecord(r)
		if err != nil {
			return nil, err
		}
		// make sure to update the meta information
		meta := pb.metas[pb.pinned]
		meta.isDirty = true
		meta.freeSpace = p.header.FreeSpace()
		// all is well, lets return
		return rid, nil
	}
	// otherwise, we must allocate enough pages and
	// link them and split the record up and write
	// the data to the pages

	return nil, nil
}

func (pb *PageBuffer) GetRecord(rid *RecordID) ([]byte, error) {
	return nil, nil
}

func (pb *PageBuffer) DelRecord(rid *RecordID) error {
	return nil
}

func (pb *PageBuffer) FreeSpace() int {
	var totalFreeSpace uint16
	for _, meta := range pb.metas {
		totalFreeSpace += meta.freeSpace
	}
	return int(totalFreeSpace)
}

func (pb *PageBuffer) DirtyPages() int {
	var dirtyPages int
	for _, meta := range pb.metas {
		if meta.isDirty {
			dirtyPages++
		}
	}
	return dirtyPages
}

func (pb *PageBuffer) Range(fn func(rid *RecordID) bool) {

}

func (pb *PageBuffer) Flush() error {
	return nil
}

func (pb *PageBuffer) Close() error {
	err := pb.manager.Close()
	if err != nil {
		return err
	}
	return nil
}
