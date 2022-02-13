package _pager

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"text/tabwriter"
)

func BtoKB(b uint64) uint64 {
	return b / 1024
}

func BtoMB(b uint64) uint64 {
	return b / 1024 / 1024
}

func BtoGB(b uint64) uint64 {
	return b / 1024 / 1024 / 1024
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage2() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v KB", BtoKB(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KB", BtoKB(m.TotalAlloc))
	fmt.Printf("\tSys = %v KB", BtoKB(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func PrintStats(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Printf("\t[MEASURMENT]\t[BYTES]\t\t[KB]\t\t[MB]\t[GC=%d]\n", mem.NumGC)
	fmt.Printf("\tmem.Alloc:\t\t%d\t%d\t\t%d\n", mem.Alloc, BtoKB(mem.Alloc), BtoMB(mem.Alloc))
	fmt.Printf("\tmem.TotalAlloc:\t%d\t%d\t\t%d\n", mem.TotalAlloc, BtoKB(mem.TotalAlloc), BtoMB(mem.TotalAlloc))
	fmt.Printf("\tmem.HeapAlloc:\t%d\t%d\t\t%d\n", mem.HeapAlloc, BtoKB(mem.HeapAlloc), BtoMB(mem.HeapAlloc))
	fmt.Printf("\t-----\n\n")
}

func PrintStatsTab(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 4, 4, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Alloc\tTotalAlloc\tHeapAlloc\tNumGC\t")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t\n", mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.NumGC)
	fmt.Fprintln(w, "-----\t-----\t-----\t-----\t")
	w.Flush()
}

func Pack2U32(dst *uint64, src1, src2 uint32) {
	*dst = uint64(src1) | uint64(src2)<<32
}

func Unpack2U32(dst *uint64) (uint32, uint32) {
	return uint32(*dst), uint32(*dst >> 32)
}

var processingPool = sync.Pool{
	New: func() interface{} {
		return []reflect.Value{}
	},
}

func isNativeType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

// Sizeof returns the estimated memory usage of object(s) not just the size of the type.
// On 64bit Sizeof("test") == 12 (8 = sizeof(StringHeader) + 4 bytes).
func Sizeof(objs ...interface{}) (sz uint64) {
	refmap := make(map[uintptr]bool)
	processing := processingPool.Get().([]reflect.Value)[:0]
	for i := range objs {
		processing = append(processing, reflect.ValueOf(objs[i]))
	}

	for len(processing) > 0 {
		var val reflect.Value
		val, processing = processing[len(processing)-1], processing[:len(processing)-1]
		if !val.IsValid() {
			continue
		}

		if val.CanAddr() {
			refmap[val.Addr().Pointer()] = true
		}

		typ := val.Type()

		sz += uint64(typ.Size())

		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				break
			}
			if refmap[val.Pointer()] {
				break
			}

			fallthrough
		case reflect.Interface:
			processing = append(processing, val.Elem())

		case reflect.Struct:
			sz -= uint64(typ.Size())
			for i := 0; i < val.NumField(); i++ {
				processing = append(processing, val.Field(i))
			}

		case reflect.Array:
			if isNativeType(typ.Elem().Kind()) {
				break
			}
			sz -= uint64(typ.Size())
			for i := 0; i < val.Len(); i++ {
				processing = append(processing, val.Index(i))
			}
		case reflect.Slice:
			el := typ.Elem()
			if isNativeType(el.Kind()) {
				sz += uint64(val.Len()) * uint64(el.Size())
				break
			}
			for i := 0; i < val.Len(); i++ {
				processing = append(processing, val.Index(i))
			}
		case reflect.Map:
			if val.IsNil() {
				break
			}
			kel, vel := typ.Key(), typ.Elem()
			if isNativeType(kel.Kind()) && isNativeType(vel.Kind()) {
				sz += uint64(kel.Size()+vel.Size()) * uint64(val.Len())
				break
			}
			keys := val.MapKeys()
			for i := 0; i < len(keys); i++ {
				processing = append(processing, keys[i])
				processing = append(processing, val.MapIndex(keys[i]))
			}
		case reflect.String:
			sz += uint64(val.Len())
		}
	}
	processingPool.Put(processing)
	return
}
