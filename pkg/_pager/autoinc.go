package _pager

import "sync"

type autoPageID struct {
	sync.Mutex
	id uint32
}

func (a *autoPageID) getNewPageID() (id uint32) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

func (a *autoPageID) undoGetNewPageID() {
	a.Lock()
	defer a.Unlock()
	a.id--
	return
}
