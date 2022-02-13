package main

import (
	"bytes"
	"fmt"
	"github.com/cagnosolutions/pager/pkg/pager"
	"github.com/cagnosolutions/pager/pkg/util"
	"runtime"
	"time"
)

func generateData(data []byte, repeat int) []byte {
	return bytes.Repeat(data, repeat)
}

func checkMemoryBefore() *runtime.MemStats {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	return mem
}

func printMemoryCheckAfter(stats *runtime.MemStats) {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	if mem.Alloc <= stats.Alloc {
		fmt.Printf("Mem allocated: 0 MB\n")
	} else {
		fmt.Printf("Mem allocated: %3.3f MB (%3.3f KB)\n", float64(mem.Alloc-stats.Alloc)/(1024*1024), float64(mem.Alloc-stats.Alloc)/(1024))
	}
}

func test() {

	fmt.Printf("[1]\n")
	util.PrintMemUsage()

	var data [][]byte
	for i := 0; i < 5; i++ {
		b := make([]byte, 0, 999999)
		data = append(data, b)
		fmt.Println(">>>", _pager.Sizeof(&data))
		fmt.Printf("[2]\n")
		util.PrintMemUsage()
		time.Sleep(1 * time.Second)
	}
	data = nil
	fmt.Printf("[3]\n")
	util.PrintMemUsage()

	runtime.GC()
	fmt.Printf("[4]\n")
	util.PrintMemUsage()
}

func main() {

	_pager.Remove("path/to/data.db")

	mem := checkMemoryBefore()

	// open an instance of the page manager
	mgr, err := _pager.OpenPageManager("path/to/data.db")
	if err != nil {
		panic(err)
	}

	printMemoryCheckAfter(mem)

	// allocate a new page using the page manager
	pg := mgr.AllocatePage()

	printMemoryCheckAfter(mem)

	_, err = pg.AddRecord(generateData([]byte{0xDE, 0xAD, 0xBE, 0xEF}, 1000))
	if err != nil {
		panic(err)
	}

	printMemoryCheckAfter(mem)

	err = mgr.WritePage(pg)
	if err != nil {
		panic(err)
	}

	printMemoryCheckAfter(mem)

	pg = mgr.AllocatePage()

	printMemoryCheckAfter(mem)

	_, err = pg.AddRecord(generateData([]byte{0xDE, 0xAD, 0xBE, 0xEF}, 1000))
	if err != nil {
		panic(err)
	}

	printMemoryCheckAfter(mem)

	err = mgr.WritePage(pg)
	if err != nil {
		panic(err)
	}

	printMemoryCheckAfter(mem)

	printMemoryCheckAfter(mem)

	mgr.Close()
}

func mainz() {

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	// open an instance of the page manager
	mgr, err := _pager.OpenPageManager("path/to/data.db")
	if err != nil {
		panic(err)
	}

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	fmt.Printf("[state] opened page manager:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// allocate a new page using the page manager
	pg := mgr.AllocatePage()

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	fmt.Printf("[state] allocated a new page:\npage.size=%d\n", _pager.Sizeof(pg))

	// add some data to the page
	id, err := pg.AddRecord(generateData([]byte{0xDE, 0xAD, 0xBE, 0xEF}, 1000))
	if err != nil {
		panic(err)
	}
	// **it should be noted that the data we wrote to
	// the page only exists in memory right now

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	fmt.Printf("[state] added ~4000 byte record to page:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can use the page manager to persist the data
	// to disk
	err = mgr.WritePage(pg)
	if err != nil {
		panic(err)
	}

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	fmt.Printf("[state] wrote page to disk:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can also bring the page back into memory
	// by using the page manager to read the record
	// back off of the disk using the PageID
	pg, err = mgr.ReadPage(id.PageID)
	if err != nil {
		panic(err)
	}

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	fmt.Printf("[state] read page from disk:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can also return a record from the page
	data, err := pg.GetRecord(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id=%v, data=%q\n", id, data)
	fmt.Printf("[state] got record from page:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	// we can also delete a record from the page
	err = pg.DelRecord(id)
	if err != nil {
		panic(err)
	}

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	// or we can delete an entire page using the page
	// manager and passing it the PageID
	err = mgr.DeletePage(id.PageID)
	if err != nil {
		panic(err)
	}
	// **it should be noted that the page manager holds
	// a reference to any pages that have been deleted,
	// so it can recycle and use them later

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)

	// and lastly, we can of course close the page manager
	err = mgr.Close()
	if err != nil {
		panic(err)
	}

	// timeout
	fmt.Println("sleeping for 30 seconds...")
	time.Sleep(30 * time.Second)
}

func runMain() {
	// open an instance of the page manager
	mgr, err := _pager.OpenPageManager("path/to/data.db")
	if err != nil {
		panic(err)
	}

	fmt.Printf("[state] opened page manager:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	var mem runtime.MemStats
	// allocate a new page using the page manager
	pg := mgr.AllocatePage()
	_pager.PrintStats(mem)

	fmt.Printf("[state] allocated a new page:\npage.size=%d\n", _pager.Sizeof(pg))

	// add some data to the page
	id, err := pg.AddRecord(generateData([]byte{0xDE, 0xAD, 0xBE, 0xEF}, 1000))
	if err != nil {
		panic(err)
	}
	// **it should be noted that the data we wrote to
	// the page only exists in memory right now

	fmt.Printf("[state] added ~4000 byte record to page:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can use the page manager to persist the data
	// to disk
	err = mgr.WritePage(pg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[state] wrote page to disk:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can also bring the page back into memory
	// by using the page manager to read the record
	// back off of the disk using the PageID
	pg, err = mgr.ReadPage(id.PageID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[state] read page from disk:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can also return a record from the page
	data, err := pg.GetRecord(id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id=%v, data=%q\n", id, data)

	fmt.Printf("[state] got record from page:\nmanager.size=%d\n", _pager.Sizeof(mgr))

	// we can also delete a record from the page
	err = pg.DelRecord(id)
	if err != nil {
		panic(err)
	}

	// or we can delete an entire page using the page
	// manager and passing it the PageID
	err = mgr.DeletePage(id.PageID)
	if err != nil {
		panic(err)
	}
	// **it should be noted that the page manager holds
	// a reference to any pages that have been deleted,
	// so it can recycle and use them later

	// and lastly, we can of course close the page manager
	err = mgr.Close()
	if err != nil {
		panic(err)
	}
}
