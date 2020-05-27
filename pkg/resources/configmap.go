package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
)

func NewXrootdConfigMapResource(xrootd *v1alpha1.Xrootd) Resource {
	configmap := objects.NewContainerConfigMap(xrootd)
	return Resource{Object: configmap}
}
