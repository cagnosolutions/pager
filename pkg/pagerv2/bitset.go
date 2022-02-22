package pagerv2

import (
	"fmt"
	"strconv"
)

const (
	ws    uint = 64                 // match arch, if 32bit, change to 32
	lg2ws uint = 6                  // log2(ws), so 6 for 64, and 5 for 32
	all   uint = 0xffffffffffffffff // aka, (1<<64)-1, 1 left shift 64-1
	max        = ^uint(0)           // max uint value for your architecture
)

// random note: lgws can also be found using bitwise operations (x >> 1)

// BitSet is a bit set data type
type BitSet struct {
	length uint
	bits   []uint
}

// alignedSize aligns a given size so it works well
func alignedSize(size uint) uint {
	if size > (max - ws + 1) {
		return max >> lg2ws
	}
	return (size + (ws - 1)) >> lg2ws
}

// resize adds additional words to incorporate new bits if needed
func (b *BitSet) resize(i uint) {
	if i < b.length || i > max {
		return
	}
	nsize := int(alignedSize(i + 1))
	if b.bits == nil {
		b.bits = make([]uint, nsize)
	} else if cap(b.bits) >= nsize {
		b.bits = b.bits[:nsize] // fast resize
	} else if len(b.bits) < nsize {
		newset := make([]uint, nsize, 2*nsize) // increase capacity 2x
		copy(newset, b.bits)
		b.bits = newset
	}
	b.length = i + 1
}

// NewBitSet sets up and returns a new *BitSet
func NewBitSet(bits uint) *BitSet {
	alignedLen := alignedSize(bits)
	return &BitSet{
		length: bits,
		bits:   make([]uint, alignedLen),
	}
}

// SetMany is identical to calling Set repeatedly
func (b *BitSet) SetMany(ii ...uint) *BitSet {
	for _, i := range ii {
		b.Set(i)
	}
	return b
}

// Set sets bit i to 1. The capacity of the bitset grows accordingly
func (b *BitSet) Set(i uint) *BitSet {
	b.resize(i) // resize if need be
	// if i >= b.length {
	//	return b
	// }
	b.bits[i>>lg2ws] |= 1 << (i & (ws - 1))
	return b
}

// Unset clears bit i, aka sets it to 0
func (b *BitSet) Unset(i uint) *BitSet {
	if i >= b.length {
		return b
	}
	b.bits[i>>lg2ws] &^= 1 << (i & (ws - 1))
	return b
}

// IsSet tests and returns a boolean if bit i is set
func (b *BitSet) IsSet(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.bits[i>>lg2ws]&(1<<(i&(ws-1))) != 0
}

// Value returns the value
func (b *BitSet) Value(i uint) uint {
	if i >= b.length {
		return 0
	}
	return b.bits[i>>lg2ws] & (1 << (i & (ws - 1)))
}

// Len returns the number of bits in the bitset
func (b *BitSet) Len() uint {
	return b.length
}

// print binary value of bitset
func (b *BitSet) String() string {
	var res = strconv.Itoa(int(ws))
	var str string
	for i := 0; i < len(b.bits); i++ {
		str += fmt.Sprintf(
			"[%d]uint%d=[%."+res+"b] (%s bits)", i, ws,
			b.bits[i], res,
		)
	}
	return str
}

func (b *BitSet) GetRawSet() []uint {
	return b.bits
}

func (b *BitSet) PercentageFull() (int, float64) {
	var isset int
	for i := uint(0); i < b.length; i++ {
		if b.bits[i>>lg2ws]&(1<<(i&(ws-1))) != 0 {
			isset++
		}
	}
	return isset, float64(isset) / float64(b.length)
}
