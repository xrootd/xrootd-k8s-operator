package resources

import (
	"github.com/shivanshs9/xrootd-operator/pkg/resources/objects"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
)

func (irs *InstanceResourceSet) AddXrootdRedirectorStatefulSetResource() {
	xrootd := irs.xrootd
	component := constant.XrootdRedirector
	objectName := utils.GetObjectName(component, xrootd.Name)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	statefulset := objects.GenerateXrootdStatefulSet(xrootd, objectName, labels, component)
	irs.addResource(Resource{Object: &statefulset})
}

func (irs *InstanceResourceSet) AddXrootdWorkerStatefulSetResource() {
	xrootd := irs.xrootd
	component := constant.XrootdWorker
	objectName := utils.GetObjectName(component, xrootd.Name)
	labels := utils.GetComponentLabels(component, xrootd.Name)
	statefulset := objects.GenerateXrootdStatefulSet(xrootd, objectName, labels, component)
	irs.addResource(Resource{Object: &statefulset})
}
