package reflects

import "reflect"

func ValueOf(i interface{}) reflect.Value {
	switch x := i.(type) {
	case reflect.Value:
		return x
	default:
		return reflect.ValueOf(i)
	}
}
