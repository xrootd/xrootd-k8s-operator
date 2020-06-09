package utils

import (
	"github.com/shivanshs9/ty/fun"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func MergeLabels(labels ...constant.Labels) constant.Labels {
	return constant.Labels(fun.MergeMaps(labels).(map[string]string))
}
