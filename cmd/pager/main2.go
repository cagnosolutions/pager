package main

import (
	"fmt"
	"github.com/cagnosolutions/pager/pkg/pager"
)

func main2() {

	// open an instance of the page manager
	manager, err := pager.OpenPageManager("path/to/data.db")
	if err != nil {
		panic(err)
	}

	// allocate a new page using the page manager
	page := manager.AllocatePage()

	// add some data to the page, notice when you write a
	// record to the page it returns a record id
	id, err := page.AddRecord([]byte("this is record one"))
	if err != nil {
		panic(err)
	}
	// **it should be noted that the data we wrote to
	// the page only exists in memory right now

	// we can use the page manager to persist the data
	// to disk
	err = manager.WritePage(page)
	if err != nil {
		panic(err)
	}

	// we can also bring the page back into memory
	// by using the page manager to read the record
	// back off of the disk using the PageID
	page, err = manager.ReadPage(id.PageID)
	if err != nil {
		panic(err)
	}

	// we can also return a record from the page
	data, err := page.GetRecord(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id=%v, data=%q\n", id, data)

	// we can also delete a record from the page
	err = page.DelRecord(id)
	if err != nil {
		panic(err)
	}

	// or we can delete an entire page using the page
	// manager and passing it the PageID
	err = manager.DeletePage(id.PageID)
	if err != nil {
		panic(err)
	}
	// **it should be noted that the page manager holds
	// a reference to any pages that have been deleted,
	// so it can recycle and use them later

	// and lastly, we can of course close the page manager
	err = manager.Close()
	if err != nil {
		panic(err)
	}
}
