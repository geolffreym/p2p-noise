package utils

import (
	"reflect"
	"testing"
)

func TestUtilClear(t *testing.T) {

	type S struct {
		a int
		b int
		c []byte
	}

	a := []byte("abcdef")
	b := []int{1, 2, 3, 4, 5, 6}
	c := S{
		a: 1,
		b: 2,
		c: a,
	}

	Clear(&a)
	Clear(&b)
	Clear(&c)

	if len(a) > 0 || len(b) > 0 {
		t.Errorf("Expected empty slices after clearing")
	}

	if (!reflect.DeepEqual(c, S{})) {
		t.Errorf("Expected empty struct after clearing")
	}

}
