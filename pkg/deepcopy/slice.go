package deepcopy

func IntSlice(s []int) []int {
	copied := make([]int, len(s))
	copy(copied, s)

	return copied
}

func InterfaceSlice(s []interface{}) []interface{} {
	copied := make([]interface{}, 0, len(s))

	for _, v := range s {
		mapVal, isMap := v.(map[string]interface{})
		if isMap {
			copied = append(copied, StringInterfaceMap(mapVal))

			continue
		}

		sliceVal, isSlice := v.([]interface{})
		if isSlice {
			copied = append(copied, InterfaceSlice(sliceVal))

			continue
		}

		copied = append(copied, v)
	}

	return copied
}
