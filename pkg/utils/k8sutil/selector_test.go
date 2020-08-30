package k8sutil

import (
	"testing"

	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestXrootdSelectorMatchesSyncedRedirectorResources(t *testing.T) {
	instance := &xrootdv1alpha1.XrootdCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
	}
	selector, err := GetXrootdLabelSelector(constant.XrootdRedirector, instance.Name)
	if err != nil {
		t.Error(err)
	}
	labels := utils.GetComponentLabels(constant.XrootdRedirector, instance.Name)
	if !selector.Matches(labels) {
		t.Errorf("selector (%v) doesn't match with label (%v)", selector, labels)
	}
}
