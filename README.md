# pager
pager is a page manager

---

Go to section
* [Hello](#hello)
* [Hello World](#hello-world)
* [Another section](#new-section) <-- it's called 'Another section' in this list but refers to 'New section'

Here is an [example label].

[example label]: #the-anchor-name-or-id "Optional title for mouse hover"

## Hello
### Hello World
## New section

## Table of Contents
- [A test](#a-test)
- [Getting Started](#getting-started)
    - [Installing](#installing)
    - [Opening a database](#opening-a-database)
    - [Transactions](#transactions)
        - [Read-write transactions](#read-write-transactions)
        - [Read-only transactions](#read-only-transactions)
        - [Batch read-write transactions](#batch-read-write-transactions)
        - [Managing transactions manually](#managing-transactions-manually)
    - [Using buckets](#using-buckets)
    - [Using key/value pairs](#using-keyvalue-pairs)
    - [Autoincrementing integer for the bucket](#autoincrementing-integer-for-the-bucket)
    - [Iterating over keys](#iterating-over-keys)
        - [Prefix scans](#prefix-scans)
        - [Range scans](#range-scans)

### A Test 

### API overview for the **page* type
```go

    // returns a new page instance (most of the
    // time you will not use this directly)
    NewPage(pid uint32) *page

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

### API overview for the **PageManager* type
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
    Close()

```
---

## Getting Started

Import the package
```go
package main

import (
	"fmt"
	"github.com/cagnosolutions/pager/pkg/pager"
)

func main() {
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