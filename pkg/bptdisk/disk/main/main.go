package main

import (
	"encoding/hex"
	"fmt"
	"github.com/cagnosolutions/pager/pkg/bptdisk/disk"
)

func main() {
	n := disk.NewNode()
	b := make([]byte, 4096)
	disk.Encode(b, n)
	fmt.Println(hex.Dump(b))
}
