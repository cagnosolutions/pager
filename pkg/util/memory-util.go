package util

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
)

func GetMemStats() *runtime.MemStats {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return m
}

func CompareMemStats(first, second *runtime.MemStats) {
	if first.Alloc > second.Alloc {
		fmt.Printf("No memory allocated!\n")
		return
	}
	fms1 := FormatMemStats("1st", first, KB)
	fms2 := FormatMemStats("2nd", second, KB)
	if len(fms1) != len(fms2) {
		fmt.Printf("Memory profiles are not in sync!\n")
		return
	}
	var sb strings.Builder
	fmt.Printf("Comparing memory stats\n=====================\n")
	for i := 0; i < len(fms1); i++ {
		sb.WriteString(fms1[i])
		sb.WriteString("\n")
		sb.WriteString(fms2[i])
		sb.WriteString("\n\n")
	}
	fmt.Print(sb.String())
	//alloc := fmt.Sprintf("Allocated: %3.3f MB\n", bToMb(second.Alloc-first.Alloc))
}

func CompareMemStatsDiff(first, second *runtime.MemStats) {
	if first.Alloc > second.Alloc {
		fmt.Printf("No memory allocated!\n")
		return
	}
	fms1 := FormatMemStats("1st", first, KB)
	fms2 := FormatMemStats("2nd", second, KB)
	if len(fms1) != len(fms2) {
		fmt.Printf("Memory profiles are not in sync!\n")
		return
	}
	var sb strings.Builder
	fmt.Printf("Comparing memory stats\n=====================\n")
	for i := 0; i < len(fms1); i++ {
		n1 := strings.Index(fms1[i], ": ")
		n2 := strings.Index(fms2[i], ": ")
		num1 := fms1[i][n1+2:]
		num2 := fms2[i][n2+2:]
		if num1 == num2 {
			continue
		}
		nt1, err := strconv.Atoi(num1)
		if err != nil {
			log.Panicf("Opps, error converting string to number: %d (%v)\n", nt1, fms1[i][:n1])
		}
		nt2, err := strconv.Atoi(num2)
		if err != nil {
			log.Panicf("Opps, error converting string to number: %d (%v)\n", nt2, fms2[i][:n2])
		}
		var diffn int
		var diffs string
		if nt1 > nt2 {
			diffn = nt1 - nt2
			diffs = "[1st] Difference: %d"
		}
		if nt2 > nt1 {
			diffn = nt2 - nt1
			diffs = "[2nd] Difference: %d"
		}
		sb.WriteString(fms1[i])
		sb.WriteString("\n")
		sb.WriteString(fms2[i])
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf(diffs, diffn))
		sb.WriteString("\n\n")
	}
	fmt.Print(sb.String())
	//alloc := fmt.Sprintf("Allocated: %3.3f MB\n", bToMb(second.Alloc-first.Alloc))
}

func PrintMemStats(m *runtime.MemStats) {
	fm := FormatMemStats("info", m, KB)
	var sb strings.Builder
	for i := 0; i < len(fm); i++ {
		sb.WriteString(fm[i])
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	fmt.Print(sb.String())
	fmt.Printf("\n\nAlloc = %v KB, %v MB", bToKb(m.Alloc), bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KB, %v MB", bToKb(m.TotalAlloc), bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KB, %v MB", bToKb(m.Sys), bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

func format(b uint64, resolution int) float32 {
	switch resolution {
	default:
		return float32(b)
	case KB:
		return bToKb(b)
	case MB:
		return bToMb(b)
	case GB:
		return bToGb(b)
	}
}
func bToGb(b uint64) float32 {
	return float32(b / 1024 / 1024 / 1024)
}

func bToMb(b uint64) float32 {
	return float32(b / 1024 / 1024)
}

func bToKb(b uint64) float32 {
	return float32(b / 1024)
}

func FormatMemStats(msg string, m *runtime.MemStats, f int) []string {
	return []string{
		fmt.Sprintf("[%v] Alloc: %v", msg, format(m.Alloc, f)),
		fmt.Sprintf("[%v] TotalAlloc: %v", msg, format(m.TotalAlloc, f)),
		fmt.Sprintf("[%v] Sys: %v", msg, format(m.Sys, f)),
		fmt.Sprintf("[%v] Lookups: %v", msg, format(m.Lookups, f)),
		fmt.Sprintf("[%v] Mallocs: %v", msg, format(m.Mallocs, f)),
		fmt.Sprintf("[%v] Frees: %v", msg, format(m.Frees, f)),
		fmt.Sprintf("[%v] HeapAlloc: %v", msg, format(m.HeapAlloc, f)),
		fmt.Sprintf("[%v] HeapSys: %v", msg, format(m.HeapSys, f)),
		fmt.Sprintf("[%v] HeapIdle: %v", msg, format(m.HeapIdle, f)),
		fmt.Sprintf("[%v] HeapInuse: %v", msg, format(m.HeapInuse, f)),
		fmt.Sprintf("[%v] HeapReleased: %v", msg, format(m.HeapReleased, f)),
		fmt.Sprintf("[%v] HeapObjects: %v", msg, format(m.HeapObjects, f)),
		fmt.Sprintf("[%v] StackInuse: %v", msg, format(m.StackInuse, f)),
		fmt.Sprintf("[%v] StackSys: %v", msg, format(m.StackSys, f)),
		fmt.Sprintf("[%v] MSpanInuse: %v", msg, format(m.MSpanInuse, f)),
		fmt.Sprintf("[%v] MSpanSys: %v", msg, format(m.MSpanSys, f)),
		fmt.Sprintf("[%v] MCacheInuse: %v", msg, format(m.MCacheInuse, f)),
		fmt.Sprintf("[%v] MCacheSys: %v", msg, format(m.MCacheSys, f)),
		fmt.Sprintf("[%v] BuckHashSys: %v", msg, format(m.BuckHashSys, f)),
		fmt.Sprintf("[%v] GCSys: %v", msg, format(m.GCSys, f)),
		fmt.Sprintf("[%v] OtherSys: %v", msg, format(m.OtherSys, f)),
		fmt.Sprintf("[%v] NextGC: %v", msg, format(m.NextGC, f)),
		fmt.Sprintf("[%v] LastGC: %v", msg, format(m.LastGC, f)),
		fmt.Sprintf("[%v] PauseTotalNs: %v", msg, format(m.PauseTotalNs, f)),
		fmt.Sprintf("[%v] NumGC: %v", msg, format(uint64(m.NumGC), f)),
		fmt.Sprintf("[%v] NumForcedGC: %v", msg, format(uint64(m.NumForcedGC), f)),
		fmt.Sprintf("[%v] GCCPUFraction: %v", msg, format(uint64(m.GCCPUFraction), f)),
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v KB, %v MB", bToKb(m.Alloc), bToMb(m.Alloc))
	fmt.Printf("\t\tTotalAlloc = %v KB, %v MB", bToKb(m.TotalAlloc), bToMb(m.TotalAlloc))
	fmt.Printf("\t\tSys = %v KB, %v MB", bToKb(m.Sys), bToMb(m.Sys))
	fmt.Printf("\t\tNumGC = %v\n", m.NumGC)
}
