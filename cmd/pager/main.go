package main

import (
	"fmt"

	"github.com/cagnosolutions/pager/pkg/_pager"
)

func _main() {
	// main2()
	// testOverflowPage()
	// pageMan()
	// testRecordPrefix()
}

func testOverflowPage() {
	// db file name
	name := fmt.Sprintf("cmd/storage/pageman/data/overflow-%.4x.db", 1)

	// remove file
	err := _pager.Remove(name)
	if err != nil {
		fmt.Errorf("[file] remove: %s", err)
	}

	// open file
	f, err := _pager.OpenPageManager(name)
	if err != nil {
		fmt.Errorf("[file] open: %s", err)
	}

	// insert data larger than one page
	p1 := f.AllocatePage()
	p2 := f.AllocatePage()
	p1 = p1.Link(p2)
	data := makeData(64, []byte{0xDE, 0xAD, 0xBE, 0xEF})

	_, err = p1.AddRecord(data[:32])
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}

	_, err = p2.AddRecord(data[32:])
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}

	f.Range(
		p1.PageID(), func(rid *_pager.RecordID) bool {
			rec, _ := p1.GetRecord(rid)
			fmt.Printf("%v => %q\n", rid, rec)
			return true
		},
	)

	// p1.Range(func(rid *pager.RecordID) bool {
	//	rec, _ := p1.GetRecord(rid)
	//	fmt.Printf("%v => %q\n", rid, rec)
	//	return true
	// })
	// p2.Range(func(rid *pager.RecordID) bool {
	//	rec, _ := p2.GetRecord(rid)
	//	fmt.Printf("%v => %q\n", rid, rec)
	//	return true
	// })

	// sort the records
	// p.SortRecords()

	// close file
	err = f.Close()
	if err != nil {
		fmt.Errorf("[file] close: %s", err)
	}
}

func makeData(dataLength int, data []byte) []byte {
	b := make([]byte, dataLength)
	for i := 0; i < dataLength; i += len(data) {
		copy(b[i:i+len(data)], data)
	}
	return b
}

func testRecordPrefix() {
	// db file name
	name := fmt.Sprintf("cmd/storage/pageman/data/file-%.4x.db", 1)

	// remove file
	err := _pager.Remove(name)
	if err != nil {
		fmt.Errorf("[file] remove: %s", err)
	}

	// open file
	f, err := _pager.OpenPageManager(name)
	if err != nil {
		fmt.Errorf("[file] open: %s", err)
	}

	// insertion order: Leslie Smith, Jonny Burkholder, Tom Wallace, Ron McKelvey
	//    sorted order: Jonny Burkholder, Leslie Smith, Ron McKelvey, Tom Wallace

	// insert some records out of lexicographic order
	p := f.AllocatePage()
	_, err = p.AddRecord([]byte("Leslie Smith"))
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}
	_, err = p.AddRecord([]byte("Jonny Burkholder"))
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}
	_, err = p.AddRecord([]byte("Tom Wallace"))
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}
	_, err = p.AddRecord([]byte("Ron McKelvey"))
	if err != nil {
		fmt.Errorf("[page] add record: %s", err)
	}

	// range the records
	p.Range(
		func(rid *_pager.RecordID) bool {
			rec, _ := p.GetRecord(rid)
			fmt.Printf("%v => %q\n", rid, rec)
			return true
		},
	)

	// sort the records
	// p.SortRecords()

	// close file
	err = f.Close()
	if err != nil {
		fmt.Errorf("[file] close: %s", err)
	}

}

func pageMan() {

	// db file name
	name := fmt.Sprintf("cmd/storage/pageman/data/file-%.4x.db", 1)

	// remove file
	err := _pager.Remove(name)
	if err != nil {
		fmt.Errorf("[file] remove: %s", err)
	}

	// open file
	f, err := _pager.OpenPageManager(name)
	if err != nil {
		fmt.Errorf("[file] open: %s", err)
	}

	// create a record id "holder"
	pgmap := make(map[uint32][]*_pager.RecordID)

	for i := 0; i < 8; i++ {
		// allocate a new page
		pg := f.AllocatePage()
		var recs []*_pager.RecordID
		// add some records to it
		for j := 0; j < 16; j++ {
			rec := fmt.Sprintf("this-is-record-number-%.6x", j)
			rid, err := pg.AddRecord([]byte(rec))
			if err != nil {
				fmt.Errorf("[page] adding record: %s", err)
			}
			recs = append(recs, rid)
		}
		// add records to page map
		pgmap[pg.PageID()] = recs
		// save page
		err = f.WritePage(pg)
		if err != nil {
			fmt.Errorf("[file] write: %s", err)
		}
	}

	// close file
	err = f.Close()
	if err != nil {
		fmt.Errorf("[file] close: %s", err)
	}

	// open file again
	f, err = _pager.OpenPageManager(name)
	if err != nil {
		fmt.Errorf("[file] open: %s", err)
	}

	// print any free pages
	pids := f.GetFreePageIDs()
	fmt.Printf("free pages found: %d\n", len(pids))
	for _, pid := range pids {
		fmt.Printf("free page: %d\n", pid)
	}

	// range the entire page map
	for pid, _ := range pgmap {
		// delete even records
		// on even page ids
		if pid%2 == 0 {
			// get page
			pg, err := f.ReadPage(pid)
			if err != nil {
				fmt.Errorf("[file] read: %s", err)
			}
			// range page
			pg.Range(
				func(rid *_pager.RecordID) bool {
					// isolate even records
					if rid.SlotID%2 == 0 {
						// remove record if it's even
						err = pg.DelRecord(rid)
						if err != nil {
							fmt.Errorf("[page] delete: %s", err)
						}
					}
					return true
				},
			)
			// save page
			err = f.WritePage(pg)
			if err != nil {
				fmt.Errorf("[file] write: %s", err)
			}
		} else {
			// and completely delete any
			// page (the entire page) if
			// it has an odd page id
			err = f.DeletePage(pid)
			if err != nil {
				fmt.Errorf("[file] delete: %s", err)
			}
		}
	}

	// see which pages are "free"
	pids = f.GetFreePageIDs()
	fmt.Printf("free pages found: %d\n", len(pids))
	for _, pid := range pids {
		fmt.Printf("free page: %d\n", pid)
	}

	// load page two
	pg2, err := f.ReadPage(2)
	if err != nil {
		fmt.Errorf("[file] read page: %s", err)
	}

	// print page
	fmt.Printf("[THIS IS PAGE 2]\n%s\n", pg2)

	// range page two
	pg2.Range(
		func(rid *_pager.RecordID) bool {
			fmt.Printf("page=2, rid=%v\n", rid)
			return true
		},
	)

	// add some records
	id, err := pg2.AddRecord([]byte("new record #2"))
	if err != nil {
		fmt.Errorf("[page 2] add record: %s", err)
	}
	fmt.Printf("added record, got id: %v\n", id)

	id, err = pg2.AddRecord([]byte("new record #4"))
	if err != nil {
		fmt.Errorf("[page 2] add record: %s", err)
	}
	fmt.Printf("added record, got id: %v\n", id)

	id, err = pg2.AddRecord([]byte("new record #6"))
	if err != nil {
		fmt.Errorf("[page 2] add record: %s", err)
	}
	fmt.Printf("added record, got id: %v\n", id)

	// print page
	fmt.Printf("[THIS IS PAGE 2]\n%s\n", pg2)

	// range page two
	pg2.Range(
		func(rid *_pager.RecordID) bool {
			fmt.Printf("page=2, rid=%v\n", rid)
			return true
		},
	)

	// save page
	err = f.WritePage(pg2)
	if err != nil {
		fmt.Errorf("[file] write: %s", err)
	}

	// close file
	err = f.Close()
	if err != nil {
		fmt.Errorf("[file] close: %s", err)
	}

}

func perr(err error) {
	if err != nil {
		fmt.Errorf("[file] remove: %s", err)
	}
}
