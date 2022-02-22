package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func Benchmark_EncodingAndDecodingAllocations(b *testing.B) {

	p := make([]byte, 64)
	pc := make([]byte, 64)
	b.ResetTimer()
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		i := 0
		// encoding uint16's into byte slice
		binary.LittleEndian.PutUint16(p[i:i+2], 11)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 12)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 13)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 14)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 21)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 22)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 23)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 24)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 31)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 32)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 33)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 34)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 41)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 42)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 43)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 44)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 51)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 52)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 53)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 54)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 61)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 62)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 63)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 64)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 71)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 72)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 73)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 74)
		i += 2
		// next...
		binary.LittleEndian.PutUint16(p[i:i+2], 81)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 82)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 83)
		i += 2
		binary.LittleEndian.PutUint16(p[i:i+2], 84)
		i += 2

		// copy the contents of p into pc
		copy(pc, p)

		i = 0
		// decoding uint64's from the byte slice
		i64b1 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b2 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b3 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b4 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b5 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b6 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b7 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8
		i64b8 := binary.LittleEndian.Uint64(p[i : i+8])
		i += 8

		// clearing the byte slice
		for i := range p {
			p[i] = 0x00
		}

		i = 0
		// encoding uint64's into byte slice
		binary.LittleEndian.PutUint64(p[i:i+8], i64b1)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b2)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b3)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b4)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b5)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b6)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b7)
		i += 8
		binary.LittleEndian.PutUint64(p[i:i+8], i64b8)
		i += 8
	}

	fmt.Printf("are p and pc the same? %v\n", bytes.Equal(p, pc))

	// fmt.Println(hex.Dump(p))
	// fmt.Printf("%x\n", p)

}

func BenchmarkNewPager1(b *testing.B) {

	b.ResetTimer()
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		p1 := NewPager1(256)
		for i := 0; i < p1.pages; i++ {
			pg := p1.GetPage(uint32(i))
			pg.SetData([]byte("this is data for a particular page"))
			pd := pg.GetData()
			_ = pd
		}
		_ = p1
	}
}

func BenchmarkNewPager2(b *testing.B) {

	b.ResetTimer()
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		p2 := NewPager2(256)
		for i := 0; i < p2.pages; i++ {
			pg := p2.GetPage(uint32(i))
			pg.SetData([]byte("this is data for a particular page"))
			pd := pg.GetData()
			_ = pd
		}
		_ = p2
	}
}
