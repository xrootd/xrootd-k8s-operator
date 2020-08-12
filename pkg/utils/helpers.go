package utils

import (
	"strings"

	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
)

// MergeLabels merges the given Labels to return one Labels object
func MergeLabels(args ...types.Labels) types.Labels {
	result := make(types.Labels, len(args)*len(args[0]))
	for _, pairs := range args {
		for key, value := range pairs {
			result[key] = value
		}
	}
	return result
}

// GetObjectName gets the object name by joining CR name and component
func GetObjectName(component types.ComponentName, crName string) types.ObjectName {
	return types.ObjectName(SuffixName(crName, string(component)))
}

// SuffixName joins the given suffixes with the given name using - as delimiter
func SuffixName(name string, suffix string, suffixes ...string) string {
	return strings.Join(append([]string{name, suffix}, suffixes...), "-")
}

// GetComponentLabels returns the Labels object for the given component
func GetComponentLabels(component types.ComponentName, crName string) types.Labels {
	labels := map[string]string{
		"component": string(component),
		"instance":  crName,
	}
	return MergeLabels(constant.ControllerLabels, labels)
}
