package pager

import (
	"fmt"
	"log"
	"testing"
)

var recIDs []uint

func addRawPageRecords(pg page) []uint {
	var recs []uint
	for i := 0; i < 10; i++ {
		rec := fmt.Sprintf("this-is-record-%.6x", i)
		rid, err := pg.AddRecord([]byte(rec))
		if err != nil {
			panic(err)
		}
		recs = append(recs, rid)
	}
	return recs
}

func TestRawPage_AddRecord(t *testing.T) {
	pg := NewRawPage(1)
	for i := 0; i < 10; i++ {
		rec := fmt.Sprintf("this-is-record-%.6x", i)
		rid, err := pg.AddRecord([]byte(rec))
		if err != nil {
			t.Errorf("[Page] adding record: %s", err)
		}
		recIDs = append(recIDs, rid)
	}
}

func TestRawPage_GetRecord(t *testing.T) {
	pg := NewRawPage(1)
	recs := addRawPageRecords(pg)
	for _, rid := range recs {
		rec, err := pg.GetRecord(rid)
		if err != nil {
			t.Errorf("[Page] getting record: %s", err)
		}
		fmt.Printf("%v: %q\n", rid, rec)
	}
}

func TestRawPage_DelRecord(t *testing.T) {
	pg := NewRawPage(1)
	log.Printf("adding records...\n")
	recs := addRawPageRecords(pg)
	log.Printf("getting records...\n")
	// get them to prove they are there
	for _, rid := range recs {
		rec, err := pg.GetRecord(rid)
		if err != nil {
			t.Errorf("[Page] getting record: %s", err)
		}
		fmt.Printf("%v: %q\n", rid, rec)
	}
	log.Printf("deleting records...\n")
	// attempt to delete half of them
	for i, rid := range recs {
		if i%2 == 0 {
			err := pg.DelRecord(rid)
			if err != nil {
				t.Errorf("[Page] deleting record: %s", err)
			}
		}
	}
	log.Printf("getting records (again)...\n")
	// get them again, to see if they are gone
	for i, rid := range recs {
		rec, err := pg.GetRecord(rid)
		if i%2 != 0 {
			if err != nil {
				t.Errorf("[Page] getting record: %s", err)
			}
			fmt.Printf("%v: %q\n", rid, rec)
		} else {
			fmt.Printf("%v: freed\n", rid)
		}
	}
}

func TestRawPage_Range(t *testing.T) {
	pg := NewRawPage(1)
	log.Println("adding records...")
	addRawPageRecords(pg)
	log.Println("ranging records...")
	pg.Range(
		func(rid uint) bool {
			rec, err := pg.GetRecord(rid)
			if err != nil {
				t.Errorf("[Page] getting record: %s", err)
			}
			fmt.Printf("%v: %q\n", rid, rec)
			return true
		},
	)
}

func Test_RawPageHeader_FreeSpace(t *testing.T) {
	pg := NewRawPage(1)
	estimatedUsedSpace := 1000 + (99 * pageSlotSize)
	estimatedFreeSpace := MaxRecordSize - estimatedUsedSpace
	for i := 0; i < 100; i++ {
		// write 10 bytes in each record
		_, err := pg.AddRecord([]byte("=========~"))
		if err != nil {
			t.Errorf("[Page] adding record: %s", err)
		}
	}
	fmt.Printf("%s\n", pg)
	freeSpace := pg.getFreeSpaceUpper() - pg.getFreeSpaceLower()
	if estimatedFreeSpace != int(freeSpace) {
		t.Errorf(
			"[Page] estimatedFreeSpace=%d, actualFreeSpace=%d",
			estimatedFreeSpace, freeSpace,
		)
	}

}
