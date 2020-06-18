package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func (irs *InstanceResourceSet) AddXrootdRedirectorConfigMapResource() {
	xrootd := irs.xrootd
	objectName := utils.GetObjectName(constant.CfgXrootd, xrootd.Name)
	labels := utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name)
	etcConfigmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.CfgXrootd, "etc")
	runConfigmap := objects.GenerateContainerConfigMap(xrootd, objectName, labels, constant.CfgXrootd, "run")
	irs.addResource(Resource{Object: &etcConfigmap}, Resource{Object: &runConfigmap})
}
