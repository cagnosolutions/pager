package main

import (
	"fmt"
	"github.com/cagnosolutions/pager/pkg/_pager"
	"github.com/cagnosolutions/pager/pkg/util"
	"runtime"
	"time"
)

func main() {
	runPageIdea()
}

// 64 KB object
type largeObject struct {
	data [^uint16(0)]byte
}

func mapAllocations() {

	// create map
	m := make(map[int]*largeObject, 1)

	// force GC, then run memstats
	//runtime.GC()

	// print memory usage
	util.PrintMemUsage()

	for i := 0; i < 10; i++ {

		// Allocate memory
		fmt.Println("adding another large object")
		m[i] = new(largeObject)

		// print memory usage
		util.PrintMemUsage()
		time.Sleep(1 * time.Second)
	}

	// "free" half of our memory
	for i := 0; i < 5; i++ {
		// removing a large object
		fmt.Println("removing a large object")
		m[i] = nil
		delete(m, i)
		// print memory usage
		util.PrintMemUsage()
		time.Sleep(1 * time.Second)
	}

	// forcing GC
	fmt.Println("forcing GC")
	runtime.GC()

	// print memory usage
	util.PrintMemUsage()

}

func pageAllocations() {
	// memory should be no higher than 150KB
	util.PrintMemUsage()

	var overall []*_pager.Page
	for i := 0; i < 256; i++ {

		// Allocate memory
		p := _pager.NewPage(uint32(i))
		overall = append(overall, p)

		// Print our memory usage at each interval
		util.PrintMemUsage()
		time.Sleep(1 * time.Second)
	}

	// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
	//overall = nil
	overall = nil
	util.PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	//runtime.GC()
	util.PrintMemUsage()
}
