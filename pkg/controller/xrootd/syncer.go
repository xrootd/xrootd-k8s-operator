package xrootd

import (
	"reflect"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/RHsyseng/operator-utils/pkg/resource/read"
	"github.com/RHsyseng/operator-utils/pkg/resource/write"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/resources"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/comparator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileXrootd) SyncResources(instance controllerutil.Object) error {
	xrootd := instance.(*xrootdv1alpha1.Xrootd)
	log := log.WithName("syncResources")

	irs := resources.NewInstanceResourceSet(xrootd)
	irs.AddXrootdRedirectorStatefulSetResource()
	irs.AddXrootdConfigMapResource()
	irs.AddXrootdRedirectorServiceResource()
	irs.AddXrootdWorkerServiceResource()
	irs.AddXrootdWorkerStatefulSetResource()
	deployed, err := r.getDeployedResources(xrootd)
	if err != nil {
		return err
	}
	writer := write.New(r.GetClient()).WithOwnerController(xrootd, r.GetScheme())
	deltaMap := comparator.GetComparator().Compare(deployed, irs.GetResources().GetK8SResources())
	for resType, delta := range deltaMap {
		if !delta.HasChanges() {
			log.Info("No changes detected")
		}
		log.Info("Processing delta", "create", len(delta.Added), "update", len(delta.Updated), "delete", len(delta.Removed), "type", resType)
		added, err := writer.AddResources(delta.Added)
		if err != nil {
			return err
		}
		updated, err := writer.UpdateResources(deployed[resType], delta.Updated)
		if err != nil {
			return err
		}
		removed, err := writer.RemoveResources(delta.Removed)
		if err != nil {
			return err
		}
		if added || updated || removed {
			log.Info("Executed changes", "added", added, "updated", updated, "removed", removed)
		}
	}
	// lockedresources, err := irs.ToLockedResources()
	// err = r.UpdateLockedResources(xrootd, lockedresources, []lockedpatch.LockedPatch{})
	// if err != nil {
	// 	log.Error(err, "unable to update locked resources", "locked resources", lockedresources)
	// }
	return nil
}

func (r *ReconcileXrootd) GetOwnedResourceKinds(instance runtime.Object) []runtime.Object {
	return []runtime.Object{
		&corev1.ConfigMapList{},
		&appsv1.StatefulSetList{},
		&corev1.ServiceList{},
	}
}

func (r *ReconcileXrootd) getDeployedResources(xrootd *xrootdv1alpha1.Xrootd) (map[reflect.Type][]resource.KubernetesResource, error) {
	reader := read.New(r.GetClient()).WithNamespace(xrootd.Namespace).WithOwnerObject(xrootd)
	return reader.ListAll(r.GetOwnedResourceKinds(xrootd)...)
}
