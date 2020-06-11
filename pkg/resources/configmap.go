package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func NewXrootdRedirectorConfigMapResource(xrootd *v1alpha1.Xrootd) Resource {
	objectName := utils.GetObjectName(constant.XrootdRedirector, constant.ConfigMap)
	labels := utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name)
	configmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.Xrootd, "etc")
	return Resource{Object: &configmap}
}
