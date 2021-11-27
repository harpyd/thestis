package deepcopy

func IntSlice(s []int) []int {
	copied := make([]int, len(s))
	copy(copied, s)

	return copied
}
