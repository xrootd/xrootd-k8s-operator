package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func (irs *InstanceResourceSet) AddXrootdRedirectorServiceResource() {
	xrootd := irs.xrootd
	component := constant.XrootdRedirector
	objectName := utils.GetObjectName(component, xrootd.Name)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	service := objects.GenerateXrootdRedirectorService(xrootd, objectName, labels, component)
	irs.addResource(Resource{Object: &service})
}
