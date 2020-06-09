package utils

import (
	"reflect"
	"testing"

	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func AssertDeep(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", expected, actual)
	}
}

func TestMergeLabels(t *testing.T) {
	l1 := constant.Labels(map[string]string{"app": "l1"})
	l2 := constant.Labels(map[string]string{"instance": "l2"})
	AssertDeep(t, constant.Labels(map[string]string{"app": "l1", "instance": "l2"}), MergeLabels(l1, l2))
}
