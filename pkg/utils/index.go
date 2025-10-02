package utils

// TransformIndex converts the provided 1-based index to a 0-based index while
// clamping the value to the valid boundaries of the slice. If the slice is
// empty, the function always returns 0.
func TransformIndex[T any](arr []T, index int) int {
	if len(arr) == 0 {
		return 0
	}

	idx := index - 1
	return min(max(idx, 0), len(arr)-1)
}
