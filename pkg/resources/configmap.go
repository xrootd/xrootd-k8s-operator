package resources

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/resources/objects"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
)

func (irs *InstanceResourceSet) AddXrootdConfigMapResource() {
	xrootd := irs.xrootd
	objectName := utils.GetObjectName(constant.CfgXrootd, xrootd.Name)
	labels := utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name)
	etcConfigmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.CfgXrootd, "etc")
	runConfigmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.CfgXrootd, "run")
	irs.addResource(Resource{Object: &etcConfigmap}, Resource{Object: &runConfigmap})
}
