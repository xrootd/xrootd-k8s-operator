package template

import (
	"reflect"
	"testing"
)

func AssertDeep(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", expected, actual)
	}
}

func TestIterateCountForPositiveCount(t *testing.T) {
	const count = 2
	AssertDeep(t, []int{0, 1}, iterateCount(count))
}

func TestIterateCountForZeroCount(t *testing.T) {
	const count = 0
	AssertDeep(t, []int{}, iterateCount(count))
}

func TestIterateCountForNegativeCount(t *testing.T) {
	const count = -1

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should have panicked!")
		} else if r != nil && r.(error).Error() != "runtime error: makeslice: len out of range" {
			t.Errorf("wrong error: %v", r)
		}
	}()

	iterateCount(count)

	t.Errorf("did not panic")
}
