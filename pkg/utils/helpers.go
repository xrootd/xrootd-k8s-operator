package utils

import (
	"strings"

	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	. "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
)

func MergeLabels(args ...Labels) Labels {
	result := make(Labels, len(args)*len(args[0]))
	for _, pairs := range args {
		for key, value := range pairs {
			result[key] = value
		}
	}
	return result
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
