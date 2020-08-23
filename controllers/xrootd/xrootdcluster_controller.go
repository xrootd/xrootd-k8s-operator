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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
)

// XrootdClusterReconciler reconciles a XrootdCluster object
type XrootdClusterReconciler struct {
	Log logr.Logger
	reconciler.BaseReconciler
	*reconciler.WatchManager
}

// blank assignment to verify that ReconcileXrootd implements reconcile.Reconciler
var _ reconcile.Reconciler = &XrootdClusterReconciler{}

// blank assignment to verify that ReconcileXrootd implements reconciler.SyncReconciler
var _ reconciler.SyncReconciler = &XrootdClusterReconciler{}

// blank assignment to verify that ReconcileXrootd implements reconciler.WatchReconciler
var _ reconciler.WatchReconciler = &XrootdClusterReconciler{}

// blank assignment to verify that ReconcileXrootd implements reconciler.StatusReconciler
var _ reconciler.StatusReconciler = &XrootdClusterReconciler{}

// +kubebuilder:rbac:groups=xrootd.xrootd.org,resources=xrootdclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=xrootd.xrootd.org,resources=xrootdclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=catalog.xrootd.org,resources=xrootdversions,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods;services;services/finalizers;endpoints;persistentvolumeclaims;events;configmaps;secrets,verbs=create;delete;get;list;patch;update;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=create;delete;get;list;patch;update;watch

// Reconcile executes the reconciliation logic on trigger of watched events
func (r *XrootdClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// ctx := context.Background()
	logger := r.Log.WithValues("xrootdcluster", req.NamespacedName)
	logger.Info("Reconciling Xrootd...")

	// Fetch the Xrootd instance
	instance := &xrootdv1alpha1.XrootdCluster{}
	// Call common reconcile logic
	result, err := reconciler.Reconcile(r, req, instance, logger)
	if err == nil {
		logger.Info("Reconciled successfully!")
	}
	return result, err
}

// IsValid determines if a Xrootd instance is valid and initializes empty fields.
func (r *XrootdClusterReconciler) IsValid(instance controllerutil.Object) (result bool, err error) {
	result = true
	xrootd := instance.(*xrootdv1alpha1.XrootdCluster)
	if xrootd.Spec.Redirector.Replicas == 0 {
		xrootd.Spec.Redirector.Replicas = 1
	}
	if xrootd.Spec.Worker.Replicas == 0 {
		xrootd.Spec.Worker.Replicas = 1
	}
	if len(xrootd.Spec.Worker.Storage.Class) == 0 {
		xrootd.Spec.Worker.Storage.Class = "standard"
	}
	if len(xrootd.Spec.Version) == 0 {
		result, err = false, fmt.Errorf("Provide xrootd version in instance")
	} else if versionInfo, tErr := utils.GetXrootdVersionInfo(r.GetClient(), instance.GetNamespace(), xrootd.Spec.Version); tErr != nil {
		result, err = false, errors.Wrapf(tErr, "Unable to find requested version - %s", xrootd.Spec.Version)
	} else if image := versionInfo.Spec.Image; len(image) == 0 {
		result, err = false, fmt.Errorf("Invalid image, '%s', provided for the given version, '%s'", image, xrootd.Spec.Version)
	} else {
		xrootd.SetVersionInfo(*versionInfo)
	}

	if result {
		xrootd.Status.SetSpecValidCondition(true, "Spec is valid", "'IsValid' check passed")
	} else {
		// if invalid spec, update invalid cluster state
		xrootd.Status.SetSpecValidCondition(false, err.Error(), "'IsValid' check failed")
	}
	return
}

// ManageError handles any error during reconciliation and updates CR status phase and condition
func (r *XrootdClusterReconciler) ManageError(instance controllerutil.Object, err error, log logr.Logger) (reconcile.Result, error) {
	xrootd := instance.(*xrootdv1alpha1.XrootdCluster)

	// set member status to empty array (otherwise status update fails)
	xrootd.Status.RedirectorStatus = xrootdv1alpha1.NewMemberStatus([]string{}, []string{})
	xrootd.Status.WorkerStatus = xrootdv1alpha1.NewMemberStatus([]string{}, []string{})

	// set cluster to failed status
	xrootd.Status.SetPhase(xrootdv1alpha1.ClusterPhaseFailed)
	if tErr := r.GetClient().Status().Update(context.Background(), xrootd); tErr != nil {
		r.Log.Error(tErr, "failed updating xrootd instance status")
		err = errors.Wrap(err, tErr.Error())
	}
	return r.BaseReconciler.ManageError(instance, err, log)
}

// SetupWithManager assigns controller manager and watches
func (r *XrootdClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.AddXrootdLogger()
	if err := r.StartWatching(); err != nil {
		return errors.Wrap(err, "failed starting watches")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&xrootdv1alpha1.XrootdCluster{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

// ControllerName is the name of xrootd controller
const ControllerName string = constant.ControllerName
