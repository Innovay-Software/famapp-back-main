package utils

func Int64PointerToUint64Pointer(item *int64) *uint64 {
	var ans *uint64
	if item != nil {
		t := uint64(*item)
		ans = &t
	}
	return ans
}