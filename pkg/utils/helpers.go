package utils

import (
	"strings"

	"github.com/shivanshs9/ty/fun"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	. "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"k8s.io/apimachinery/pkg/labels"
)

func MergeLabels(args ...Labels) Labels {
	return fun.MergeMaps(args).(labels.Set)
}

func GetObjectName(component ComponentName, controllerName string) ObjectName {
	return ObjectName(SuffixName(controllerName, string(component)))
}

func SuffixName(name string, suffix string, suffixes ...string) string {
	return strings.Join(append([]string{name, suffix}, suffixes...), "-")
}

func GetComponentLabels(component ComponentName, controllerName string) Labels {
	labels := map[string]string{
		"component": string(component),
		"instance":  controllerName,
	}
	return MergeLabels(constant.ControllerLabels, labels)
}
