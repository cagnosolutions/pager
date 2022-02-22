package tests

import (
	"encoding/binary"
)

const pageSize = 4 << 10

type Pager1 struct {
	data  []*Page1
	pages int
}

func NewPager1(pages int) *Pager1 {
	p := &Pager1{
		data:  make([]*Page1, pages),
		pages: pages,
	}
	for i := 0; i < pages; i++ {
		p.data[i] = NewPage1(uint32(i), 0, 0)
	}
	return p
}

func (p *Pager1) GetPage(pid uint32) *Page1 {
	return p.data[pid]
}

func (p *Page1) SetData(rec []byte) {
	n := copy(p.Data, rec)
	p.Size = uint32(n)
}

func (p *Page1) GetData() []byte {
	return p.Data[:p.Size]
}

type Page1 struct {
	PageID uint32
	NextID uint32
	PrevID uint32
	Size   uint32
	Data   []byte
}

func NewPage1(pid, prev, next uint32) *Page1 {
	return &Page1{
		PageID: pid,
		NextID: 0,
		PrevID: 0,
		Size:   0,
		Data:   make([]byte, pageSize),
	}
}

type Pager2 struct {
	data  []byte
	pages int
}

//go:noescape
func NewPager2(pages int) *Pager2 {
	p := &Pager2{
		data:  make([]byte, pages*pageSize),
		pages: pages,
	}
	return p
}

type Page2 []byte

func NewPage2(p Page2, pid, prev, next uint32) {
	binary.LittleEndian.PutUint32(p[0:4], pid)
	binary.LittleEndian.PutUint32(p[4:8], prev)
	binary.LittleEndian.PutUint32(p[8:12], next)
}

func (p *Pager2) GetPage(pid uint32) Page2 {
	off := pid * pageSize
	return p.data[off : off+pageSize]
}

func (p Page2) SetData(rec []byte) {
	n := copy(p[16:], rec)
	binary.LittleEndian.PutUint32(p[12:16], uint32(n))
}

func (p Page2) GetData() []byte {
	size := binary.LittleEndian.Uint32(p[12:16])
	return p[16 : 16+size]
}
