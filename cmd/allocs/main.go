package main

import (
	"github.com/cagnosolutions/pager/pkg/pager"
	"github.com/cagnosolutions/pager/pkg/util"
	"time"
)

func main() {
	// memory should be no higher than 150KB
	util.PrintMemUsage()

	var overall []*pager.Page
	for i := 0; i < 256; i++ {

		// Allocate memory
		p := pager.NewPage(uint32(i))
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
