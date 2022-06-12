package noise

import (
	"fmt"
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

func TestIndexOf(t *testing.T) {
	slice := []int{1, 2, 3, 4, 6, 7, 8, 9}
	// Table driven test
	// For each expected event
	for _, e := range slice {
		t.Run(fmt.Sprintf("Match for %x", e), func(t *testing.T) {
			match := IndexOf(slice, e)
			if ^match == 0 {
				t.Errorf("expected matched existing elements in slice index: %#v", e)
			}
		})

	}
}

func TestInvalidIndexOf(t *testing.T) {
	slice := []int{1, 2, 3, 4, 6, 7, 8, 9}
	match := IndexOf(slice, 5)

	if ^match != 0 {
		t.Error("Number 5 is not in slice cannot be found")
	}

}
