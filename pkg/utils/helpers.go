package utils

import (
	"fmt"

	"github.com/shivanshs9/ty/fun"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	. "github.com/shivanshs9/xrootd-operator/pkg/utils/types"
)

func MergeLabels(labels ...Labels) Labels {
	return Labels(fun.MergeMaps(labels).(map[string]string))
}

func GetObjectName(component ComponentName, kind KindName) ObjectName {
	return ObjectName(SuffixName(string(kind), string(component)))
}

func SuffixName(name string, suffix string) string {
	return fmt.Sprintf("%s-%s", name, suffix)
}

func GetComponentLabels(component ComponentName, controllerName string) Labels {
	labels := map[string]string{
		"component": string(component),
		"instance":  controllerName,
	}
	return MergeLabels(constant.ControllerLabels, labels)
}
