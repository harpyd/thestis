package deepcopy

func IntSlice(s []int) []int {
	copied := make([]int, 0, len(s))
	copy(copied, s)

	return copied
}
