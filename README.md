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
implemented using a slotted-page structure. This allows for easier sorting
of records within a page and also makes working with variable sized records 
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
  - [Using the manager](#using-the-page-manager)
  - [Using pages](#using-pages)
  - [Using records](#using-records)
  - [Swapping pages](#swapping-pages)
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

### Using the page manager
The `*PageManager` is one of the main top-level objects in this package. 
It is represented as a single file on your disk. To use the manager use
the `OpenPageManager(path)` method of the pager package.
```go
mgr, err := pager.OpenPageManager("path/data.db")
if err != nil {
    panic(err)
}
defer m.Close()
```

To close the manager, use the manager's `Close()` method.
```go
// to close the manager
err = manager.Close()
if err != nil {
	panic(err)
}
```
There are a lot of other operations that are available through the manager
relating to page management such as `AllocatePage()`, `DeletePage(pid)`, 
`WritePage(page)` and `ReadPage(pid)`. These will be covered below under
the sub categories [using pages](#using-pages) and 
[swapping pages](#swapping-pages).

### Using pages
Each page holds collection of one or more records. Each page has a page 
ID, which is represented as a `uint32`. To allocate a new page, use the
manager's `AllocatePage()` method.
```go
// allocate a new page
pg := mgr.AllocatePage()
```
Each page must have a unique ID. The `*PageManager` enforces unique page 
ID's when allocating. *You can also allocate a new page directly by calling
`NewPage(id uint32)`, but this is not recommended. If you do choose to use
this second method of page allocation, it is expected that you will manage 
and enforce the uniqueness of the page IDs.*

The page's ID is an important thing to keep track of. To get an allocated
page's ID, use the page's `PageID()` method.

```go
// get the page id
pid := pg.PageID()
```

To delete a page, use the manager's `DeletePage(pid)` method and provide it
with a valid page ID. It will return an error if the provided ID is invalid. 
```go
// delete page with pid
err = mgr.DeletePage(pid)
if err != nil {
	panic(err)
}
// It should be noted that the page 
// manager holds a reference to any 
// pages that have been deleted, so 
// it can recycle and use them later
```

### Using records
A record is just an arbitrary slice of bytes. To add a record to a page, use 
the `AddRecord(data)` method of the page. *Records are prefix sorted by the 
page (using the first 8 bytes of the record data) in lexicographic order.*
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
to a page, it returns a `*RecordID` type which can be used to retrieve it,
or delete it. The action of adding a new record may produce an error if the
record is empty, larger than the page itself, or if the page is out of room.

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

Use the `DelRecord(id)` method to delete a record from the page. It only
returns an error if the provided `*RecordID` is not valid.
```go
// delete record using id
err = pg.DelRecord(id)
if err != nil {
    panic(err)
}
```

### Swapping pages
***The term "swapping", in the context of this package, is used to refer to 
the action of reading and writing pages between main memory and the disk.****

Use the manager's `WritePage(page)` method to ensure any record operations or
page mutations are persisted to the disk. It will return an error if the
manager, or the underlying file has been closed. 
```go
// write page to disk
err = mgr.WritePage(pg)
if err != nil {
    panic(err)
}
```
To read a persisted page off the disk back into main memory use the manager's
`ReadPage(pid)` method. It will require a valid PageID, which can be found by
calling the page's `PageID()` method. *It can also be found by accessing the
first part of the `*RecordID` tuple of a record that resides within the page.*
```go
// use out saved pid
pid := pg.PageID()
...
// read an on disk page back into memory
pg, err = mgr.ReadPage(pid)
if err != nil {
    panic(err)
}
```

## Types

### RecordID [*][1]
A `*RecordID` is returned after a record is written to a page. It is the 
primary ID used to locate a record, return a record and to delete a record
within a page.
```go
// RecordID the unique ID of a record. It
// is a tuple containing the page ID where
// it resides, along with the slot where 
// the metadata for the record is stored.
type RecordID struct {
	PageID uint32
	SlotID uint16
}
```

### Page [*][2]
A `page` is an unexported type, and it is something that is never directly
accessed. However, it does have several methods that may be used for managing 
records. It also contains a few "getters" and "setters" for accessing a page's
header and meta-data, and an iterator for accessing a page's records in 
sequentially sorted order.
```go
// page is a page in memory and is unexported
type page struct {
    // contains unexported fields...
}

// returns a new page instance (most of the
// time you will not use this directly)
NewPage(pid uint32) *page
```
A `page` is a single contiguous block of bytes 8KB in size. It exists in 
memory only unless it is manually persisted. Any action that modifies record
data on a page is not persisted unless an explicit call to `WritePage(page)` 
is made by the manager.

### Page Manager [*][3]
```go
// PageManager manages instances of
// pages, and swapping between disk
// and memory
type PageManager struct {
    // contains unexported fields...
}

// a *PageManager takes a path in order to persist a file
func NewPageManager(path string) (*PageManager, error)
```
The `PageManager` is responsible for allocating new pages, de-allocating
pages,and managing the swapping of the pages between disk and memory. It 
also keeps a small cache of previously used pages in an attempt to reuse 
pages that have been de-allocated (aka deleted). When allocating new pages
the `PageManager` automatically generates and assigns each page with a 
unique ID.

### Page Buffer [*][4]
```go
// PageBuffer provides buffered page
// management, and access to larger
// spans of data
type PageBuffer struct {
    // contains unexported fields...
}

// a *PageBuffer instance takes a *PageManager
func NewPageBuffer(pm *PageManager) (*PageBuffer, error)
```
A `PageBuffer` wraps a `PageManager` instance and provides a buffered set of
pages (default 8 pages) to work with. One advantage to using a `PageBuffer`
is that it enables you to write records that would normally be too large to 
fit inside one page--it takes care of the inter-page linking for you. 

[1]: /pkg/pager/page.go#L13 ("record id source")
[2]: /pkg/pager/page.go#L97 ("page source")
[3]: /pkg/pager/manager.go#L10 ("page manager source")
[4]: /pkg/pager/buffer.go#L8 ("page buffer source")