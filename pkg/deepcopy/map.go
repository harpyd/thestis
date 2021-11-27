package deepcopy

func StringInterfaceMap(m map[string]interface{}) map[string]interface{} {
	copied := make(map[string]interface{}, len(m))

	for k, v := range m {
		cv, ok := v.(map[string]interface{})
		if ok {
			copied[k] = StringInterfaceMap(cv)
		} else {
			copied[k] = v
		}
	}

	return copied
}
