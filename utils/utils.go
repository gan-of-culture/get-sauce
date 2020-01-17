package utils

import "reflect"

func GetLastItem(list interface{}) interface{} {
	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Slice {
		return nil
	}
	return v.Index(v.Len() - 1)
}
