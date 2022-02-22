package main

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var x, y []int

const size = 65536

func initData() {
	x = make([]int, size)
	y = make([]int, size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// fill x and y with random data
	fill(r, x, 0.5)
	fill(r, y, 0.3)
}

func resetData() {
	x = nil
	y = nil
	runtime.GC()
}

func fill(r *rand.Rand, b []int, probability float32) {
	for i := 0; i < len(b); i++ {
		for j := uint64(0); j < 64; j++ {
			if r.Float32() < probability {
				b[i] |= 1 << j
			}
		}
	}
}

func doit(a, b, c []int) {
	// simple tight for loop, but the
	// compiler can't inline it by default
	for i := 0; i < len(a); i++ {
		c[i] = a[i] + b[i]
	}
}

func Benchmark_LoopNormal(b *testing.B) {
	initData()
	z := make([]int, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loopNormal(x, y, z)
	}
	b.StopTimer()
	resetData()
}

func loopNormal(a, b, c []int) {
	// simple tight for loop, but the
	// compiler can't inline it by default
	for i := 0; i < len(a); i++ {
		c[i] = a[i] + b[i]
	}
}

func Benchmark_LoopCompilerHint(b *testing.B) {
	initData()
	z := make([]int, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loopCompilerHint(x, y, z)
	}
	b.StopTimer()
	resetData()
}

func loopCompilerHint(a, b, c []int) {
	if len(a) != len(b) || len(b) != len(c) {
		return
	}
	// simple tight for loop, but the
	// compiler can't inline it by default
	for i := 0; i < len(a); i++ {
		c[i] = a[i] + b[i]
	}
}
