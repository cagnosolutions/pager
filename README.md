# pager
pager is a page manager

## Table of Contents
- [Getting Started](#getting-started)
- [Types](#types)
    - [RecordID](#recordid)
    - [Page](#page)
    - [Page Manager](#page-manager)
    - [Page Buffer](#page-buffer)

## Getting Started
Import the package

```go
package main

import (
	"fmt"
	"github.com/cagnosolutions/pager/pkg/pager"
)
```
Opening an instance of the ***PageManager***
```go
manager, err := pager.OpenPageManager("path/to/data.db")
if err != nil {
	panic(err)
}
```

Allocate a new page using the ***PageManager***
```go
page := manager.AllocatePage()
```
Add some data to the page--notice when you write a record to the page it returns a ***RecordID***
```go
id, err := page.AddRecord([]byte("this is record one"))
if err != nil {
	panic(err)
}
// It should be noted that the data we wrote to the page 
// only exists in memory right now
```
Use the ***PageManager*** to persist any modified page data to disk
```go
err = manager.WritePage(page)
if err != nil {
	panic(err)
}
```
Use the ***PageManager*** to read a persisted page off disk into memory (using a valid ***PageID***)
```go
page, err = manager.ReadPage(id.PageID)
if err != nil {
	panic(err)
}
```
We can also get a record from the page (using a valid ***RecordID***)
```go
data, err := page.GetRecord(id)
if err != nil {
	panic(err)
}
fmt.Printf("id=%v, data=%q\n", id, data)
```
Or delete a record from a page (using a valid ***RecordID***)
```go
err = page.DelRecord(id)
if err != nil {
	panic(err)
}
```
We can also delete an entire page using the ***PageManager*** (using a valid ***PageID***)
```go
err = manager.DeletePage(id.PageID)
if err != nil {
	panic(err)
}
// It should be noted that the page 
// manager holds a reference to any 
// pages that have been deleted, so 
// it can recycle and use them later
```
And lastly, we can of course close the ***PageManager***
```go
err = manager.Close()
if err != nil {
	panic(err)
}
```

## Types

### RecordID
```go
// represents a record id
type RecordID struct {
	PageID uint32
	SlotID uint16
}
```
a record id bla bla...

### Page
```go
// page is a page in memory and is unexported
type page struct {
	// contains unexported fields...
}

// returns a new page instance (most of the
// time you will not use this directly)
NewPage(pid uint32) *page
```
*page is a single contiguous block of bytes 8KB in size. A page exists 
in memory only unless it is persisted using the ****PageManager***. Remember,
any action that modifies record data on a page is not persisted unless an 
explicit call to *WritePage(p \*page)* is made by the ****PageManager***.
<br><br>****page is NOT guaranteed to be safe concurrently***
### Methods of *page
```go
// returns the (un-fragmented) free space 
// remaining in the page  
FreeSpace() uint

// returns a boolean indicating if this
// page has been marked "free"
PageIsFree() bool
    
// takes two pages and links their prev and
// next pointers in their headers, so they
// can be traversed in order
LinkPages(a, b *page) *page
    
// simply a method form of the LinkPages
// function above, but otherwise provides 
// identical functionality 
Link(next *page) *page

// returns the ID of the page
PageID() uint32

// returns the prev page ID pointer
PrevID() uint32

// returns the next page ID pointer
NextID() uint32
    
// checks the incoming record to see if
// the record is too small, too large or
// if there is enough room left
CheckRecord(recordSize uint16) error

// prefix sorts the records in the page
SortRecords()

// adds a new record to the page and
// returns a unique *RecordID
AddRecord(r []byte) (*RecordID, error)

// returns a record from the page using
// the supplied *RecordID
GetRecord(rid *RecordID) ([]byte, error)

// deletes a record from the page using
// the supplied *RecordID
DelRecord(rid *RecorID) error

// ranges the records in the page providing
// a functional iterator for accessing the
// page records in prefix sorted order
Range(fn func(rid *RecordID) bool)
    
// resets the page to it's initial state
Reset()
```

## Page Manager
```go
// PageManager manages pages and disk persistence
type PageManager struct {
	// contains unexported fields...
}

// a *PageManager takes a path in order to persist a file
func NewPageManager(path string) (*PageManager, error)
```
****PageManager*** is responsible for allocating new pages and managing the direct
I/O of the pages between disk and memory. It also keeps a small cache of previously
used pages in an attempt to reuse deleted pages. When allocating new pages the
****PageManager*** automatically generates and assigns each page with a unique PageID.
<br><br>****PageManager is NOT guaranteed to be safe concurrently***

### Methods of *PageManager
```go
// allocates a new page 
AllocatePage() *page

// checks for any free pages it can use and if
// none are found it will allocate a new one  
GetFreeOrAllocate() *page
    
// attempts to read and return a page off disk
// using the provided PageID
ReadPage(pid uint32) (*page, error)
    
// attempts to write an in memory page to disk
WritePage(p *page) error

// attempts to delete a page (page is marked as
// a "free" page, so it can be recycled)
DeletePage(pid uint32) error
    
// returns any PageID's for pages that have been
// removed or are listed as "free" 
GetFreePageIDs() []uint32

// returns the total number of pages the manager
// has a reference to (including any "free" pages)
PageCount() int
    
// closes the manager (and the underlying file)
Close() error
```

## Page Buffer
```go
// PageBuffer provides buffered page management
type PageBuffer struct {
	// contains unexported fields...
}

// a *PageBuffer instance takes a *PageManager
func NewPageBuffer(pm *PageManager) (*PageBuffer, error)
```
****PageBuffer*** wraps a ****PageManager*** instance and provides a buffered
set of pages (default 8 pages) to work with. One advantage to using a ****PageBuffer***
is that it enables you to write records that would normally be too large to 
fit inside one page--it takes care of the inter-page linking for you. It is 
also synchronized, so it is safe to use concurrently. A ****PageBuffer*** may 
periodically flush its contents to disk, but if you are not sure, make suer 
you make a call to Flush(). 
<br><br>****PageBuffer IS guaranteed to be safe concurrently***
### Methods of *PageBuffer

```go
// adds a new record to the pinned page and
// returns a unique *RecordID
AddRecord(r []byte) (*RecordID, error)

// returns a record from the pinned page using
// the supplied *RecordID
GetRecord(rid *RecordID) ([]byte, error)

// deletes a record from the pinned page
// using the supplied *RecordID
DelRecord(rid *RecorID) error

// returns the free space for the page associated
// with the PageID provided
FreeSpace(pid uint32) int

// returns the total available free space 
// for the combined pages in the *PageBuffer 
TotalFreeSpace() int

// returns the number of dirty pages
DirtyPages() int

// forces any dirty pages to disk
Flush() error

// closes the *PageBuffer (and the underlying *PageManager)
Close() error
```