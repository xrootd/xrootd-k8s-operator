package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
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

func TestGetObjectName(t *testing.T) {
	name := GetObjectName(constant.XrootdRedirector, "test")
	if len(name) == 0 {
		t.Errorf("got empty object name: %v", name)
	}
}

func TestSuffixNameWithOneSuffix(t *testing.T) {
	name := "hello"
	suffix := "world"
	result := SuffixName(name, suffix)
	AssertDeep(t, fmt.Sprintf("%s-%s", name, suffix), result)
}

func TestSuffixNameWithTwoSuffixes(t *testing.T) {
	name := "hello"
	suffix := "world"
	suffix2 := "test"
	result := SuffixName(name, suffix, suffix2)
	AssertDeep(t, fmt.Sprintf("%s-%s-%s", name, suffix, suffix2), result)
}

func TestGetComponentLabels(t *testing.T) {
	name := "test-xrootd"
	labels := GetComponentLabels(constant.XrootdRedirector, name)
	AssertDeep(t, string(constant.XrootdRedirector), labels["component"])
	AssertDeep(t, name, labels["instance"])
}
