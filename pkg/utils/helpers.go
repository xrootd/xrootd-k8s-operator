package utils

import (
	"strings"

	"github.com/shivanshs9/ty/fun"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	. "github.com/shivanshs9/xrootd-operator/pkg/utils/types"
)

func MergeLabels(labels ...Labels) Labels {
	return Labels(fun.MergeMaps(labels).(map[string]string))
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
