package pagerv2

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const items = 16384

// ~19,000 ns/op
func BenchmarkFindIndexInSlice(b *testing.B) {
	// b.ReportAllocs()
	// b.ResetTimer()
	a := make([]int, items)
	a[items/2] = 1
	var res int
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(a); j++ {
			if a[j] == 1 {
				res = j
				break
			}
		}
	}
	// fmt.Printf("res=%d\n", res)
	_ = res
}

// ~8 ns/op for direct access
// ~741,889 ns/op for any range calls
func BenchmarkFindIndexInMapRange(b *testing.B) {
	// b.ReportAllocs()
	// b.ResetTimer()
	a := make(map[int]int, items)
	for i := 0; i < (items); i++ {
		if i == items/2 {
			a[i] = 1
		} else {
			a[i] = 0
		}
	}
	var res int
	for i := 0; i < b.N; i++ {
		// v, ok := a[items/2]
		// if ok && v == items/2 {
		// 	res = v
		// }
		for j := 0; j < len(a); j++ {
			v, ok := a[j]
			if ok && v == 1 {
				res = v
				break
			}
		}
	}
	// fmt.Printf("res=%d\n", res)
	_ = res
}

// ~8 ns/op for direct access
// ~741,889 ns/op for any range calls
func BenchmarkFindIndexInMapDirect(b *testing.B) {
	// b.ReportAllocs()
	// b.ResetTimer()
	a := make(map[int]int, items)
	for i := 0; i < (items); i++ {
		if i == items/2 {
			a[i] = 1
		} else {
			a[i] = 0
		}
	}
	var res int
	for i := 0; i < b.N; i++ {
		v, ok := a[items/2]
		if ok && v == items/2 {
			res = v
		}
	}
	// fmt.Printf("res=%d\n", res)
	_ = res
}

// ~400 ns/op
func BenchmarkFindIndexInBitmap(b *testing.B) {
	// b.ReportAllocs()
	// b.ResetTimer()
	a := NewBitSet(items)
	a.Set(items / 2)
	var res int
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(a.bits); j++ {
			if a.bits[j] > 0 {
				for k := j; k < j+8; k++ {
					if a.IsSet(uint(k)) {
						res = k
						break
					}
				}
			}
		}
	}
	// fmt.Printf("res=%d\n", res)
	_ = res
}

func TestBitSet_IsSet(t *testing.T) {
	bs := NewBitSet(32)
	ok := bs.IsSet(0)
	if ok {
		t.Errorf("expected=%v, got=%v\n", false, ok)
	}
	fmt.Printf("%s\n", bs)
}

var bsA, bsB *BitSet
var bitmapLength = 65536

func initData() {
	bsA = NewBitSet(uint(bitmapLength))
	bsB = NewBitSet(uint(bitmapLength))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fill(r, bsA.GetRawSet(), 0.5)
	fill(r, bsB.GetRawSet(), 0.3)
}

func _BenchmarkBitSet_Normal(b *testing.B) {
	initData()
	resBitmap := make([]uint, bitmapLength)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		andnot(bsA.GetRawSet(), bsB.GetRawSet(), resBitmap)
	}
	_ = resBitmap
}

func _BenchmarkBitSet_CompilerHint(b *testing.B) {
	initData()
	resBitmap := make([]uint, bitmapLength)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		andnotInlined(bsA.GetRawSet(), bsB.GetRawSet(), resBitmap)
	}
	_ = resBitmap
}

func fill(r *rand.Rand, b []uint, probability float32) {
	for i := 0; i < len(b); i++ {
		for j := uint64(0); j < 64; j++ {
			if r.Float32() < probability {
				b[i] |= 1 << j
			}
		}
	}
}

func indexes(a []uint64) []int {
	var res []int
	for i := 0; i < len(a); i++ {
		for j := 63; j > 0; j-- {
			if a[i]&(1<<uint64(j)) > 0 {
				res = append(res, (63-j)+(i*64))
			}
		}
	}
	return res
}

func and(a []uint, b []uint, res []uint) {
	for i := 0; i < len(a); i++ {
		res[i] = a[i] & b[i]
	}
}

func andInlined(a []uint, b []uint, res []uint) {
	i := 0

loop:
	if i < len(a) {
		res[i] = a[i] & b[i]
		i++
		goto loop
	}
}

func andnot(a []uint, b []uint, res []uint) {
	for i := 0; i < len(a); i++ {
		res[i] = a[i] & ^b[i]
	}
}

func andnotInlined(a []uint, b []uint, res []uint) {
	if len(a)/8 != len(b)/8 || len(b)/8 != len(res) {
		return
	}
	i := 0
loop:
	if i < len(a) {
		res[i] = a[i] & ^b[i]
		i++
		goto loop
	}
}
