package gotrace

import "math"

func lobit(x uint) uint {
	var res uint = 32
	for x&0xFFFFFF != 0 {
		x <<= 8
		res -= 8
	}
	for x != 0 {
		x <<= 1
		res -= 1
	}
	return res
}
func hibit(x uint) uint {
	var res uint = 0
	for x > math.MaxUint8 {
		x >>= 8
		res += 8
	}
	for x != 0 {
		x >>= 1
		res += 1
	}
	return res
}
