package deepcopy

import "reflect"

func Interface(i interface{}) interface{} {
	var (
		oldVal = reflect.ValueOf(i).Elem()
		newObj = reflect.New(reflect.TypeOf(i).Elem())
		newVal = newObj.Elem()
	)

	for i := 0; i < oldVal.NumField(); i++ {
		newValField := newVal.Field(i)
		if newValField.CanSet() {
			newValField.Set(oldVal.Field(i))
		}
	}

	return newObj.Interface()
}
