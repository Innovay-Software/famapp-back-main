package utils

// Return elements that are only in left
func SliceLeftExcludeRight[T comparable](leftSlice *[]T, rightSlice *[]T) *[]T {
	seen := make(map[T]bool)
	for _, item := range *rightSlice {
		seen[item] = true
	}
	result := make([]T, 0)
	for _, item := range *leftSlice {
		if seen[item] {
			continue
		}
		result = append(result, item)
	}
	return &result
}

// Returns elements that are in either left or right but not both
func SliceXor[T comparable](leftSlice *[]T, rightSlice *[]T) *[]T {
	seen := make(map[T]int)
	for _, item := range *leftSlice {
		seen[item] += 1
	}
	for _, item := range *rightSlice {
		seen[item] += 1
	}

	result := make([]T, 0)
	for k, v := range seen {
		if v == 1 {
			result = append(result, k)
		}
	}
	return &result
}

// Returns elements that are in both left and right
func SliceInnerJoin[T comparable](leftSlice *[]T, rightSlice *[]T) *[]T {
	seen := make(map[T]int)
	for _, item := range *leftSlice {
		seen[item] = 1
	}
	result := make([]T, 0)
	for _, item := range *rightSlice {
		if seen[item] == 1 {
			result = append(result, item)
		}
	}
	return &result
}

// Reverse elements in the original slice
func ReverseSliceInPlace[T any](slice *[]T) {
	for i, j := 0, len(*slice)-1; i < j; i, j = i+1, j-1 {
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	}
}

// Reverse elements but return a new slice
func ReverseSliceToNew[T any](slice *[]T) *[]T {
	n := len(*slice)
	newSlice := make([]T, n)
	for i := 0; i < n; i++ {
		newSlice[i] = (*slice)[n-1-i]
	}
	return &newSlice
}

func Int64SliceToUint64(slice *[]int64) *[]uint64 {
	ans := []uint64{}
	for _, item := range *slice {
		ans = append(ans, uint64(item))
	}
	return &ans
}

func Uint64SliceToInt64(slice *[]uint64) *[]int64 {
	ans := []int64{}
	for _, item := range *slice {
		ans = append(ans, int64(item))
	}
	return &ans
}