package _pager

func BinarySearch(set []int, find int) int {
	// declare for later
	i, j := 0, len(set)
	// otherwise, perform binary search
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if find > set[h] {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
