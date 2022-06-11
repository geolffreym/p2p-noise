// Package utils implements shared logic for packages
package noise

import (
	"reflect"
)

// Clear garbage collectable struct/*
func Clear(v any) {
	// https://stackoverflow.com/questions/29168905/how-to-clear-values-of-a-instance-of-a-type-struct-dynamically/51006888#51006888
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

// TODO write test
// IndexOf find index for element in slice
// It return index if found else -1
func IndexOf[T comparable](collection []T, el T) int {
	for i, v := range collection {
		if v == el {
			return i
		}
	}

	return -1
}
