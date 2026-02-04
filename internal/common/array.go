package common

import (
	"reflect"
)

// flatten walks arbitrary nesting of slices/arrays and collects ints.
func Flatten(input interface{}) []int {
	v := reflect.ValueOf(input)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		var res []int
		for i := 0; i < v.Len(); i++ {
			res = append(res, Flatten(v.Index(i).Interface())...)
		}
		return res
	case reflect.Int:
		return []int{int(v.Int())}
	default:
		return []int{} // or panic / error if unexpected type
	}
}
