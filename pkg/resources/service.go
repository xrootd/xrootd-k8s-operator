package resources

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/resources/objects"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
)

// AddXrootdRedirectorServiceResource adds service for redirector component
func (irs *InstanceResourceSet) AddXrootdRedirectorServiceResource() {
	xrootd := irs.xrootd
	component := constant.XrootdRedirector
	objectName := utils.GetObjectName(component, xrootd.Name)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	service := objects.GenerateXrootdService(xrootd, objectName, labels, component)
	irs.addResource(Resource{Object: &service})
}

// AddXrootdWorkerServiceResource adds service for worker component
func (irs *InstanceResourceSet) AddXrootdWorkerServiceResource() {
	xrootd := irs.xrootd
	component := constant.XrootdWorker
	objectName := utils.GetObjectName(component, xrootd.Name)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	service := objects.GenerateXrootdService(xrootd, objectName, labels, component)
	irs.addResource(Resource{Object: &service})
}
