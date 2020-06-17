package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func (irs *InstanceResourceSet) AddXrootdRedirectorConfigMapResource() {
	xrootd := irs.xrootd
	objectName := utils.GetObjectName(constant.XrootdRedirector, xrootd.Name)
	labels := utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name)
	configmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.Xrootd, "etc")
	irs.addResource(Resource{Object: &configmap})
}
