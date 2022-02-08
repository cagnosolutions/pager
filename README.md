## Overview

***Even though pager aims to make certain lower-level operations easier, the
aim is for it to be wrapped in larger structures. Most of the package is not 
guaranteed to be thread-safe. To work with data in multiple goroutines it is
recommended that locking is used to ensure only one goroutine can have access 
at a time.**** 

Pager attempts to bring some lower-level memory, filesystem, and data 
management abstractions into a more approachable package. The aim is to keep 
the API as simple as possible while still providing the necessary functions 
required for record, page and data file management.

The pager package provides methods for easily allocating and de-allocating 
pages in memory. Simple read and write methods for swapping pages to and 
from disk. In-page record management provides easy methods for reading, 
writing, deleting and sorting records within a page. It also has a page
buffer which provides a buffered page pool along with simple methods that
wrap more complex and fine-grained operations like record overflowing and
inter-page linking.

## Technical Notes
On the more technical end of things, the pages within this package are 
implemented using a slotted-page scheme. This allows for easier sorting of 
records within a page and also makes working with variable sized records 
much easier. Each page also contains a header with previous and next page
pointers (by default, they are not filled out) which makes it fairly simple
to link pages together and create complex on disk structures like lists and
trees. 

*It should be noted that some components of the package obtain a file 
lock on the data file being used, so multiple processes cannot open the same
data file at the same time. Opening an already open data file may have 
unintended side effects, and at the very least may cause it to hang until 
the other process closes it.**

## Requirements

Requires Go version 1.17.x or later, no external dependencies are required.

## Versioning
This package uses [semantic versioning](http://semver.org), so the API may 
change between major releases but not for patches and minor releases.

## Table of Contents

- [Overview](#overview)
- [Technical Notes](#overview)
- [Requirements](#requirements)
- [Versioning](#versioning)
- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Importing](#importing)
  - [Using the manager](#opening-the-page-manager)
  - [Using pages](#using-pages)
  - [Using records](#using-records)
  - [Types](#types)
    - [RecordID](#recordid)
    - [Page](#page)
    - [Page Manager](#page-manager)
    - [Page Buffer](#page-buffer)

## Getting Started

### Installing
First, make sure you meet the [requirements](#requirements), then simply
run `go get`:
```shell
$ go get github.com/cagnosolutions/pager 
```

### Importing
To use the pager package, import as:
```go
import "github.com/cagnosolutions/pager/pkg/pager"
```

### Opening the page manager
The `*PageManager` is one of the main top-level objects in this package. 
It is represented as a single file on your disk. 
```go
mgr, err := pager.OpenPageManager("path/data.db")
if err != nil {
    panic(err)
}
defer m.Close()
```

### Using pages
Each page holds collection of one or more records. Each page has a page 
ID, which is represented as a `uint32`.
```go
// allocate a new page
pg := mgr.AllocatePage()
```
Each page must have a unique ID. The `*PageManager` enforces unique page 
ID's when allocating. *You can also allocate a new page directly by calling
`NewPage(id uint32)`, but this is not recommended. If you do choose to use
this second method of page allocation, it is expected that you will manage 
and enforce the uniqueness of the page IDs.* 

### Using records
A record is just an arbitrary slice of bytes. To add a record to a page, use 
the `AddRecord(data)` method of the page.
```go
// add a record to the page, keep track of the id
id, err := pg.AddRecord([]byte("my first record"))
if err != nil {
	panic(err)
}
// It should be noted that the data we wrote to 
// the page only exists in main memory right now
```
This will add `"my first record"` to the page `pg`. When you add a record
to a page, it returns a `*RecordID` you can use to retrieve it. The action
of adding a new record may produce an error if the record is empty, larger
than the page itself, or if the page is out of room.

To retrieve this record from the page, we can use the `GetRecord(id)` method
of the page.
```go
// get a record from a page using the id
r, err := pg.GetRecord(id)
if err != nil {
	panic(err)
}
fmt.Printf("id=%v, record=%q\n", id, r)
// output:
// id=&{1 2}, record="this is record one"
```
The `GetRecord(id)` method will return an error if the record cannot be found
which may happen if the provided `*RecordID` is not valid. 

Use the `DelRecord(id)` method to delete a record from the page.
```go
// delete record using id
err = page.DelRecord(id)
if err != nil {
	panic(err)
}
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