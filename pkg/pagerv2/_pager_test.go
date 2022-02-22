package pagerv2

import (
	"fmt"
	"testing"
)

func TestNewPageCache(t *testing.T) {

	fp, err := OpenFile("pager_test/data-001.db")
	if err != nil {
		t.Error(err)
	}
	pc := NewPageCache(fp, 64)

	for i := 0; i < 64; i++ {
		pg := pc.NewPage()
		copy(pg, "this is a test, this is only a test")
	}

	for i := 0; i < 64; i++ {
		pg, err := pc.ReadPage(uint32(i))
		if err != nil {
			t.Error(">>>", err)
		}
		fmt.Printf("page[%d]=%q\n", i, pg)
	}

	fp.Close()
}
