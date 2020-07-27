package k8sutil

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func GetXrootdLabelSelector(component types.ComponentName, instanceName string) (labels.Selector, error) {
	selector := labels.NewSelector()
	for key, value := range utils.GetComponentLabels(component, instanceName) {
		req, err := labels.NewRequirement(key, selection.Equals, []string{value})
		if err != nil {
			return nil, err
		}
		selector = selector.Add(*req)
	}
	return selector, nil
}
