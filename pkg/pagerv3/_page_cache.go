package pagerv3

import (
	"errors"
)

const (
	maxCachedPages = 8
)

type PageCache struct {
	diskManager DiskManager
	pages       [maxCachedPages]*Page
	// replacer    *ClockReplacer
	// freeList    []FrameID
	// pageTable   map[PageID]FrameID
}

// clock replacement algo
// fetch page
//

// https://github.com/brunocalza/buffer-pool-manager/blob/main/buffer_pool_manager.go

// NewPage allocates a new page in the page buffer pool with
// help from the disk manager
func (pc *PageCache) NewPage() Page {
	frameID, isFromFreeList := pc.getFrameID()
	if frameID == nil {
		return nil
	}

	if !isFromFreeList {
		// remove page from current frame
		currentPage := pc.pages[*frameID]
		if currentPage != nil {
			if currentPage.isDirty {
				pc.diskManager.WritePage(currentPage)
			}

			delete(pc.pageTable, currentPage.id)
		}
	}

	// allocates new page
	pageID := pc.diskManager.AllocatePage()
	if pageID == nil {
		return nil
	}
	page := newPage(*pageID, 1, false)

	pc.pageTable[*pageID] = *frameID
	pc.pages[*frameID] = page

	return page
}

// FetchPage fetches the requested page from the page buffer
func (pc *PageCache) FetchPage(pid PageID) Page {
	// 1a) if it is in the buffer return it
	// 2a) if it not in the buffer check for free page
	// 3a) if there are no free pages, evict a page and
	//     read the requested page from the disk
	// if it is on buffer pool return it
	if frameID, ok := pc.pageTable[pageID]; ok {
		page := pc.pages[frameID]
		page.pinCount++
		(*pc.replacer).Pin(frameID)
		return page
	}

	// get the id from free list or from replacer
	frameID, isFromFreeList := pc.getFrameID()
	if frameID == nil {
		return nil
	}

	if !isFromFreeList {
		// remove page from current frame
		currentPage := pc.pages[*frameID]
		if currentPage != nil {
			if currentPage.isDirty {
				pc.diskManager.WritePage(currentPage)
			}

			delete(pc.pageTable, currentPage.id)
		}
	}

	page, err := pc.diskManager.ReadPage(pageID)
	if err != nil {
		return nil
	}
	(*page).pinCount = 1
	pc.pageTable[pageID] = *frameID
	pc.pages[*frameID] = page

	return page
}

// UnpinPage unpins the target page from the page buffer
func (pc *PageCache) UnpinPage(pid PageID, isDirty bool) error {
	if frameID, ok := pc.pageTable[pageID]; ok {
		page := pc.pages[frameID]
		page.DecPinCount()

		if page.pinCount <= 0 {
			(*pc.replacer).Unpin(frameID)
		}

		if page.isDirty || isDirty {
			page.isDirty = true
		} else {
			page.isDirty = false
		}

		return nil
	}

	return errors.New("Could not find page")
}

// FlushPage flushes the target page to disk
func (pc *PageCache) FlushPage(pid PageID, isDirty bool) error {
	if frameID, ok := pc.pageTable[pageID]; ok {
		page := pc.pages[frameID]
		page.DecPinCount()

		pc.diskManager.WritePage(page)
		page.isDirty = false

		return true
	}

	return false
}
