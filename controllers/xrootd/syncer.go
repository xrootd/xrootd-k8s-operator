/*


Copyright (C) 2020  The XRootD Collaboration

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
USA
*/

package controllers

import (
	"reflect"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"github.com/RHsyseng/operator-utils/pkg/resource/read"
	"github.com/RHsyseng/operator-utils/pkg/resource/write"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/resources"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/comparator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *XrootdClusterReconciler) SyncResources(instance controllerutil.Object) error {
	xrootd := instance.(*xrootdv1alpha1.XrootdCluster)
	log := r.Log.WithName("syncResources")

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
	requested := irs.GetResources().GetK8SResources()
	requestedMap := compare.NewMapBuilder().Add(requested...).ResourceMap()
	deltaMap := comparator.GetComparator().Compare(deployed, requestedMap)
	for resType, delta := range deltaMap {
		logger := log.WithValues("type", resType)
		if !delta.HasChanges() {
			logger.Info("No changes detected")
			continue
		}
		logger.Info("Processing delta", "create", len(delta.Added), "update", len(delta.Updated), "delete", len(delta.Removed))
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
			logger.Info("Executed changes", "added", added, "updated", updated, "removed", removed)
		}
	}
	// lockedresources, err := irs.ToLockedResources()
	// err = r.UpdateLockedResources(xrootd, lockedresources, []lockedpatch.LockedPatch{})
	// if err != nil {
	// 	log.Error(err, "unable to update locked resources", "locked resources", lockedresources)
	// }
	return nil
}

func (r *XrootdClusterReconciler) GetOwnedResourceKinds(instance runtime.Object) []runtime.Object {
	return []runtime.Object{
		&corev1.ConfigMapList{},
		&appsv1.StatefulSetList{},
		&corev1.ServiceList{},
	}
}

func (r *XrootdClusterReconciler) getDeployedResources(xrootd *xrootdv1alpha1.XrootdCluster) (map[reflect.Type][]resource.KubernetesResource, error) {
	reader := read.New(r.GetClient()).WithNamespace(xrootd.Namespace).WithOwnerObject(xrootd)
	return reader.ListAll(r.GetOwnedResourceKinds(xrootd)...)
}
