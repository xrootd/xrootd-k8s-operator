package xrootd

import (
	"fmt"

	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_xrootd")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Xrootd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileXrootd{
		BaseReconciler: reconciler.NewBaseReconciler(
			mgr.GetClient(), mgr.GetScheme(), mgr.GetEventRecorderFor(controllerName), mgr.GetConfig()),
		WatchManager: reconciler.NewWatchManager(nil),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Xrootd
	if err = c.Watch(
		&source.Kind{Type: &xrootdv1alpha1.Xrootd{}},
		&handler.EnqueueRequestForObject{},
		predicate.GenerationChangedPredicate{},
	); err != nil {
		return err
	}

	if reconciler, ok := r.(*ReconcileXrootd); ok {
		reconciler.AddXrootdLogger()
		reconciler.StartWatching()
	}

	return nil
}

// blank assignment to verify that ReconcileXrootd implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileXrootd{}

// blank assignment to verify that ReconcileXrootd implements reconciler.SyncReconciler
var _ reconciler.SyncReconciler = &ReconcileXrootd{}

// blank assignment to verify that ReconcileXrootd implements reconciler.WatchReconciler
var _ reconciler.WatchReconciler = &ReconcileXrootd{}

// blank assignment to verify that ReconcileXrootd implements reconciler.StatusReconciler
var _ reconciler.StatusReconciler = &ReconcileXrootd{}

// ReconcileXrootd reconciles a Xrootd object
type ReconcileXrootd struct {
	reconciler.BaseReconciler
	*reconciler.WatchManager
}

// Reconcile reads that state of the cluster for a Xrootd object and makes changes based on the state read
// and what is in the Xrootd.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileXrootd) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Xrootd")

	// Fetch the Xrootd instance
	instance := &xrootdv1alpha1.Xrootd{}
	result, err := reconciler.Reconcile(r, request, instance, reqLogger)
	if err == nil {
		reqLogger.Info("Reconciled successfully!")
	}
	return result, err
}

// IsValid determines if a Xrootd instance is valid and initializes empty fields.
func (r *ReconcileXrootd) IsValid(instance controllerutil.Object) (bool, error) {
	xrootd := instance.(*xrootdv1alpha1.Xrootd)
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
		return false, fmt.Errorf("Provide xrootd version in instance")
	}
	if versionInfo, err := utils.GetXrootdVersionInfo(r.GetClient(), instance.GetNamespace(), xrootd.Spec.Version); err != nil {
		return false, errors.Wrapf(err, "Unable to find requested version - %s", xrootd.Spec.Version)
	} else if image := versionInfo.Spec.Image; len(image) == 0 {
		return false, fmt.Errorf("Invalid image, '%s', provided for the given version, '%s'", image, xrootd.Spec.Version)
	} else {
		xrootd.SetVersionInfo(*versionInfo)
	}
	return true, nil
}

const controllerName = constant.ControllerName
