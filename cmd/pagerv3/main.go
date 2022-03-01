package main

import (
	"fmt"
	"log"
)

func main() {
	p := NewPager(OpenFile(32))
	p.Write([]byte("r1,p1"), p.GetFreePID())
	p.Write([]byte("r1,p2"), p.GetFreePID())
	p.Write([]byte("r1,p3"), p.GetFreePID())
	p.Write([]byte("r1,p4"), p.GetFreePID())
	p.Write([]byte("r1,p5"), p.GetFreePID())
	p.Write([]byte("r1,p6"), p.GetFreePID())
	fmt.Printf("%s\n", p)
}

const (
	pgsz = 8
	npgs = 8
)

func align(n int, size int) int {
	return (n + size) &^ size
}

type file struct {
	d []byte
	c int
}

func OpenFile(size int) *file {
	size = align(size, pgsz)
	return &file{
		d: make([]byte, size),
		c: 0,
	}
}

type pg struct {
	beg, end int
}

type pager struct {
	f *file
	d [pgsz * npgs]byte
	e int
	i map[int]pg
	n int
}

func (p *pager) String() string {
	var ss string
	for i := 0; i < p.n; i++ {
		ss += fmt.Sprintf("pg[%d]={%s}\n", i, p.d[i*pgsz:(i+1)*pgsz])
	}
	return ss
}

func NewPager(f *file) *pager {
	p := &pager{
		f: f,
		i: make(map[int]pg, npgs),
		n: npgs,
	}
	for i := 0; i < npgs; i++ {
		p.i[i] = pg{i * pgsz, (i + 1) * pgsz}
	}
	log.Printf("%+v\n", p.i)
	return p
}

func (p *pager) pageHasRoom(pid int) bool {
	return false
}

func (p *pager) GetFreePID() int {
	if p.e+1+pgsz > len(p.d) {
		return -1
	}
	return p.e
}

func (p *pager) GetFreeRID() int {
	return -1
}

func (p *pager) Read(pid int) []byte {
	at, found := p.i[pid]
	if !found {
		return nil
	}
	page := make([]byte, pgsz)
	copy(page, p.d[at.beg:at.end])
	return page
}

func (p *pager) Write(d []byte, pid int) {
	if len(d) > pgsz {
		panic("data too large")
	}
	at, found := p.i[pid]
	if !found {
		panic("bad page id")
	}
	log.Println(at)
	copy(p.d[at.beg:at.end], d)
	p.e++
}

func (p *pager) ReadRecord(pid, rid int) []byte {
	return nil
}

func (p *pager) WriteRecord(d []byte, pid int) {
	return
}
