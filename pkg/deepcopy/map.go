package deepcopy

func StringInterfaceMap(m map[string]interface{}) map[string]interface{} {
	copied := make(map[string]interface{}, len(m))

	for k, v := range m {
		mapVal, isMap := v.(map[string]interface{})
		if isMap {
			copied[k] = StringInterfaceMap(mapVal)

			continue
		}

		sliceVal, isSlice := v.([]interface{})
		if isSlice {
			copied[k] = InterfaceSlice(sliceVal)

			continue
		}

		copied[k] = v
	}

	return copied
}
