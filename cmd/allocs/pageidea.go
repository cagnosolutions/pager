package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

func runPageIdea() {
	// open a "pager" with a 16 page buffer
	p := OpenPager(16)
	// lookup page #2
	pg := p.Lookup(2)
	// print page data (BEFORE WE TRY TO MODIFY)
	fmt.Printf("Page.Data=%q\n", pg.String())
	// modify contents
	n := copy(pg.Data[:], []byte("can we modify page data content..."))
	copy(pg.Data[n:], []byte("...without writing it back to the pager?"))
	// print page data again (AFTER WE TRY TO MODIFY)
	fmt.Printf("Page.Data=%q\n", pg.String())
}

const pageSize = 4096

type PageHeader struct {
	// empty for now
}

type Page struct {
	Header *PageHeader
	Data   [pageSize]byte
}

func (pg *Page) String() string {
	return string(pg.Data[:64]) + "..."
}

type Pager struct {
	cache map[uint32]*Page
}

func OpenPager(pages int) *Pager {
	m := make(map[uint32]*Page, pages)
	for i := 1; i <= pages; i++ {
		m[uint32(i)] = &Page{
			// empty structs and nil pointers occupy zero bytes
			Header: &PageHeader{},
			Data:   [pageSize]byte{},
		}
	}
	return &Pager{
		cache: m,
	}
}

func (p *Pager) Lookup(id uint32) *Page {
	pg, ok := p.cache[id]
	if !ok {
		return nil
	}
	return pg
}

func (p *Pager) Addr(id uint32) string {
	return fmt.Sprintf("%p", p.cache[id])
}

func (p *Pager) Size() int {
	// get initial map ptr size
	size := 16 // 16 bytes for the map ptr
	for i := 1; i <= len(p.cache); i++ {
		// add size of key and value
		// 4 bytes for the uint32 key, and 8 for
		// the page ptr, plus the actual page size
		size += 4 + 8 + int(unsafe.Sizeof(*(p.cache[uint32(i)])))
	}
	return size
}

func (p *Pager) String() string {
	data := map[string]interface{}{
		"Cache": p.cache,
		"Size":  p.Size(),
	}
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}
