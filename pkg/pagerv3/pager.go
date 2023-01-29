package pagerv3

// Pager represents the main API of a page buffer or cache.
// Individual implementations may differ quite a bit.
type Pager interface {

	// WritePage takes the data provided and writes it to the
	// page using the PageID provided. If successful a nil
	// error is returned.
	WritePage(d []byte, pid PageID) error

	// ReadPage attempts to read the contents of the page using
	// the PageID provided. If successful a nil error is returned.
	ReadPage(pid PageID) ([]byte, error)

	// SyncPage performs a file system fsync which forces the OS
	// and disk controller to flush the buffered contents of the
	// page found using the PageID provided onto the physical media.
	// If successful, a nil error is returned.
	SyncPage(pid PageID) error

	// FreePage attempts to free the page found using the PageID
	// provided. If the page is currently dirty, it first attempts
	// to call fsync. A boolean indicating the success of the free
	// call will be returned along with a nil error, if successful.
	FreePage(pid PageID) (bool, error)

	// WriteRecord takes the record data provided and writes it to
	// the page using the PageID provided. If successful a RecordID
	// will be returned along with a nil error. The RecordID can be
	// used at a later time to delete or update the specific record.
	WriteRecord(d []byte, pid PageID) (RecordID, error)

	// ReadRecord attempts to read and return a copy of the contents
	// of the selected record using the PageID and RecordID provided.
	// If successful, a nil error will be returned.
	ReadRecord(pid PageID, rid RecordID) ([]byte, error)

	// DeleteRecord attempts to delete the contents of the selected
	// record using the PageID and RecordID provided. A boolean will
	// be returned indicating the success of the call to delete and
	// if successful, a nil error will be returned.
	DeleteRecord(pid PageID, rid RecordID) error
}

type (
	PageID   uint32
	RecordID uint32
	FrameID  uint32
)

type Manager interface {

	// NewPage allocates a new page and pins it to a frame. If we did not
	// find an open frame we will proceed by attempting to victimize the
	// current frame.
	NewPage() *Page

	// FetchPage fetches the requested page from the buffer pool. If the
	// page is in cache, it is returned immediately. If not, it will be
	// found on disk, loaded into the cache and returned.
	FetchPage(pid PageID) *Page

	// UnpinPage unpins the target page from the buffer pool. It indicates
	// that the page is not used any more for the current requesting thread.
	// If no more threads are using this page, the page is considered for
	// eviction (victim).
	UnpinPage(pid PageID) error

	// FlushPage flushes the target page that is in the cache onto the
	// underlying medium. It also decrements the pin count and unpins it
	// from the holding frame and unsets the dirty bit.
	FlushPage(pid PageID) bool

	// DeletePage deletes a page from the buffer pool. Once removed, it
	// marks the holding fame as free to use.
	DeletePage(pid PageID) error

	// GetFrameID returns a frame ID from the free list, or by using the
	// replacement policy if the free list is full along with a boolean
	// indicating true if the frame ID was returned using the free list
	// and false if it was returned by using the replacement policy.
	GetFrameID() (FrameID, bool)
}
