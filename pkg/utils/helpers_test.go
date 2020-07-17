package utils

import (
	"reflect"
	"testing"

	. "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
)

func AssertDeep(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", expected, actual)
	}
}

func TestMergeLabels(t *testing.T) {
	l1 := Labels(map[string]string{"app": "l1"})
	l2 := Labels(map[string]string{"instance": "l2"})
	AssertDeep(t, Labels(map[string]string{"app": "l1", "instance": "l2"}), MergeLabels(l1, l2))
}
