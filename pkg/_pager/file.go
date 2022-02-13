package _pager

import (
	"sync"
)

//	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//	NOTE: this file will most likely be removed soon...
//	!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// File is mainly just a synchronized
// abstraction on top of a PageManager
type File struct {
	mu     sync.RWMutex
	pm     *PageManager
	pp     *Page // ptr to the most recently allocated Page
	doSync bool
}

// OpenFile opens a File
func _OpenFile(path string, sync bool) (*File, error) {
	pm, err := OpenPageManager(path)
	if err != nil {
		return nil, err
	}
	f := &File{
		pm:     pm,
		doSync: sync,
	}
	err = f.load()
	if err != nil {
		return nil, err
	}
	return f, nil
}

// load is the File's init method
func (f *File) load() error {
	if f.pm.PageCount() == 0 {
		f.pp = f.pm.AllocatePage()
		return nil
	}
	var size uint16
	var pgid uint32
	for _, h := range f.pm.pageHeaders {
		if h.FreeSpace() > size {
			size = h.FreeSpace()
			pgid = h.pageID
		}
	}
	pg, err := f.pm.ReadPage(pgid)
	if err != nil {
		return err
	}
	f.pp = pg
	return nil
}

// Write attempts to write data with as little
// hassle as possible
func (f *File) Write(record []byte) (*RecordID, error) {
	// lock
	f.mu.Lock()
	defer f.mu.Unlock()
tryagain:
	// write record to Page data
	rid, err := f.pp.AddRecord(record)
	if err != nil {
		// check to see if it's a space issue
		if err == ErrNoMoreRoomInPage {
			f.pp = f.pm.AllocatePage()
			goto tryagain
		}
		return nil, err
	}
	if f.doSync {
		err = f.pm.WritePage(f.pp)
		if err != nil {
			return rid, err
		}
	}
	return rid, nil
}

func (f *File) Read(recordID *RecordID) ([]byte, error) {
	// read lock
	f.mu.RLock()
	defer f.mu.RUnlock()
	// if the record is in the current Page
	// then we don't have to read it!
	if recordID.PageID == f.pp.header.pageID {
		// get the record from our Page cache
		rec, err := f.pp.GetRecord(recordID)
		if err != nil {
			return nil, err
		}
		// go it!
		return rec, nil
	}
	// otherwise, we must read the Page in...
	pg, err := f.pm.ReadPage(recordID.PageID)
	if err != nil {
		return nil, err
	}
	// and then read the record
	rec, err := pg.GetRecord(recordID)
	if err != nil {
		return nil, err
	}
	// got it!
	return rec, nil
}

func (f *File) Range(fn func(rid *RecordID) bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, h := range f.pm.pageHeaders {
		if h.PageIsFree() {
			continue
		}
		pg, err := f.pm.ReadPage(h.pageID)
		if err != nil {
			return err
		}
		pg.Range(func(rid *RecordID) bool {
			return fn(rid)
		})
	}
	return nil
}

func (f *File) Save() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.save()
}

func (f *File) save() error {
	err := f.pm.WritePage(f.pp)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	err := f.save()
	if err != nil {
		return err
	}
	err = f.pm.Close()
	if err != nil {
		return err
	}
	return nil
}
