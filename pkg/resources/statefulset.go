package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func NewXrootdRedirectorStatefulSetResource(xrootd *v1alpha1.Xrootd) Resource {
	component := constant.XrootdRedirector
	objectName := utils.GetObjectName(component, constant.StatefulSet)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	statefulset := objects.GenerateXrootdStatefulSet(xrootd, objectName, labels, component)
	return Resource{Object: &statefulset}
}
