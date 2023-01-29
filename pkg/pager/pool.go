package pager

import (
	"encoding/binary"
	"sync"
)

type slot struct {
	slotID uint16
	status uint16
	offset uint16
	length uint16
	prefix string
}

type slotSet []slot

var slotSetPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pageSlotSize)
	},
}

func GetSlotSet() slotSet {
	return slotSetPool.Get().(slotSet)
}

func PutSlotSetPool(ss slotSet) {

}

// a single entry is 8 bytes
type entries []byte

func (e entries) _Swap(i, j int) {
	// basically what im looking
	// to get is something like...
	// e[i:i+8], e[j:j+8] = e[j:j+8], e[i:i+8]
	// ^ but that doesn't work, lol
}

// long version
func (p page) swap(i, j int) {
	// get slot `i` and slot `j`
	slotI := p.getWholePageSlot(i)
	slotJ := p.getWholePageSlot(j)
	// encode slot `i` and slot `j` into uint64's
	slotIu64 := binary.LittleEndian.Uint64(slotI)
	slotJu64 := binary.LittleEndian.Uint64(slotJ)
	// decode uint64s `i` and `j` back into bytes
	// but cross them up thus doing our "swap"
	binary.LittleEndian.PutUint64(slotI, slotJu64)
	binary.LittleEndian.PutUint64(slotJ, slotIu64)
	// note: haven't tested this yet
}

var bin = binary.LittleEndian

// short version
func (p page) swapShort(i, j int) {
	si, sj := p.getWholePageSlot(i), p.getWholePageSlot(j)
	bin.PutUint64(si, bin.Uint64(sj))
	bin.PutUint64(sj, bin.Uint64(si))
	// note: haven't tested this yet
}
