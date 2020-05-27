package xrootd

import (
	"context"

	xrootdv1alpha1 "github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/resources"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	oputil "github.com/redhat-cop/operator-utils/pkg/util"
	lockcontroller "github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller"
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
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
		EnforcingReconciler: lockcontroller.NewEnforcingReconciler(
			mgr.GetClient(), mgr.GetScheme(), mgr.GetConfig(), mgr.GetEventRecorderFor(controllerName),
		),
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
	if err = c.Watch(&source.Kind{Type: &xrootdv1alpha1.Xrootd{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	if enforcingReconciler, ok := r.(*lockcontroller.EnforcingReconciler); ok {
		if err = c.Watch(
			&source.Channel{Source: enforcingReconciler.GetStatusChangeChannel()},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	} else {
		log.V(1).Info("The given reconciler is not EnforcingReconciler", "reconciler", r)
		return nil
	}

	return nil
}

// blank assignment to verify that ReconcileXrootd implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileXrootd{}

// ReconcileXrootd reconciles a Xrootd object
type ReconcileXrootd struct {
	lockcontroller.EnforcingReconciler
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
	err := r.GetClient().Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if ok, err := r.IsValid(instance); !ok {
		// If CR isn't valid, update the status with error and return
		return r.ManageError(instance, err)
	}

	if ok := r.IsInitialized(instance); !ok {
		// If CR isn't initialized, update the status with error and return
		err := r.GetClient().Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, "Unable to update instance", "instance", instance)
			return r.ManageError(instance, err)
		}
		return reconcile.Result{}, nil
	}

	if oputil.IsBeingDeleted(instance) {
		if !oputil.HasFinalizer(instance, controllerName) {
			return reconcile.Result{}, nil
		}
		// TODO: Write Cleanup Logic
		oputil.RemoveFinalizer(instance, controllerName)
		err = r.GetClient().Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, "unable to update instance", "instance", instance)
			return r.ManageError(instance, err)
		}
		return reconcile.Result{}, nil
	}

	if err = r.syncResources(instance); err != nil {
		return r.ManageError(instance, err)
	}

	return r.ManageSuccess(instance)
}

func (r *ReconcileXrootd) syncResources(xrootd *xrootdv1alpha1.Xrootd) error {
	resourcesList := []resources.Resource{
		resources.NewXrootdConfigMapResource(xrootd),
	}
	transformer := func(res resources.Resource) lockedresource.LockedResource {
		result, err := res.ToLockedResource()
		if err != nil {
			// TODO: how to return error in outer function scope?
		}
		return *result
	}
	lockedresources := utils.Map(transformer, resourcesList).([]lockedresource.LockedResource)
	r.UpdateLockedResources(xrootd, lockedresources)
	return nil
}

const controllerName = "xrootd-controller"
