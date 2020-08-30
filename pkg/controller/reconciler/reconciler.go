package reconciler

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type resource = controllerutil.Object

// Reconciler provides common helper methods useful for reconciliation.
type Reconciler interface {
	GetClient() client.Client
	GetRecorder() record.EventRecorder
	GetScheme() *runtime.Scheme
	GetConfig() *rest.Config
	GetResourceInstance(request reconcile.Request, instance resource) error
	ManageError(instance resource, err error, log logr.Logger) (reconcile.Result, error)
	ManageSuccess(instance resource, log logr.Logger) (reconcile.Result, error)
	IsValid(instance resource) (bool, error)
}

// BaseReconciler implements common logic for Reconciler
type BaseReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	config   *rest.Config
	recorder record.EventRecorder
}

// NewBaseReconciler creates a new BaseReconciler
func NewBaseReconciler(
	client client.Client, scheme *runtime.Scheme,
	recorder record.EventRecorder, config *rest.Config,
) BaseReconciler {
	return BaseReconciler{
		client:   client,
		scheme:   scheme,
		recorder: recorder,
		config:   config,
	}
}

// GetClient returns controller runtime client object
func (br BaseReconciler) GetClient() client.Client {
	return br.client
}

// GetRecorder implements Reconciler.GetRecorder
func (br BaseReconciler) GetRecorder() record.EventRecorder {
	return br.recorder
}

// GetScheme implements Reconciler.GetScheme
func (br BaseReconciler) GetScheme() *runtime.Scheme {
	return br.scheme
}

// GetConfig implements Reconciler.GetConfig
func (br BaseReconciler) GetConfig() *rest.Config {
	return br.config
}

// IsBeingDeleted returns whether this object has been requested to be deleted
func IsBeingDeleted(obj resource) bool {
	return !obj.GetDeletionTimestamp().IsZero() || obj.GetName() == ""
}

// GetResourceInstance implements Reconciler.GetResourceInstance by fetching the requested k8s
// object and updating the `instance` struct
func (br *BaseReconciler) GetResourceInstance(request reconcile.Request, instance resource) error {
	if err := br.GetClient().Get(context.TODO(), request.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return nil
		}
		// Error reading the object - requeue the request.
		return err
	}
	return nil
}

// ManageError manage error sets an error status in the CR and fires an event, finally it returns the error so the operator can re-attempt
func (br *BaseReconciler) ManageError(instance resource, err error, log logr.Logger) (reconcile.Result, error) {
	br.GetRecorder().Event(instance, "Warning", "ProcessingError", err.Error())
	return reconcile.Result{}, err
}

// ManageSuccess will update the status of the CR and return a successful reconcile result
func (br *BaseReconciler) ManageSuccess(instance resource, log logr.Logger) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

// IsValid signifies whether given instance is valid
func (br *BaseReconciler) IsValid(instance resource) (bool, error) {
	return true, nil
}

// Reconcile is a helper method abstracting the CR fetching, validation, syncing owned resources and watching resources
func Reconcile(r Reconciler, request reconcile.Request, instance resource, log logr.Logger) (result reconcile.Result, err error) {
	defer func() {
		// handle panics and send the error result back to event processor
		if tErr, ok := recover().(error); ok {
			if err != nil {
				err = errors.Wrap(tErr, err.Error())
			} else {
				err = errors.Wrap(tErr, "faced panic on reconciliation")
			}
		}
		if err != nil {
			result, err = r.ManageError(instance, err, log)
		}
	}()

	if err = r.GetResourceInstance(request, instance); err != nil {
		result, err = r.ManageError(instance, err, log)
		return
	}

	if IsBeingDeleted(instance) {
		log.Info("Deleting instance...", "instance", instance)
		// TODO: Write Cleanup Logic
		result, err = reconcile.Result{}, nil
		return
	}

	if ok, tErr := r.IsValid(instance); !ok {
		// If CR isn't valid, update the status with error and return
		result, err = r.ManageError(instance, tErr, log)
		return
	}

	if syncer, ok := r.(SyncReconciler); ok {
		log.Info("Started syncing resources...")
		if tErr := syncer.SyncResources(instance); tErr != nil {
			log.Error(err, "Failed syncing resources...")
			result, err = r.ManageError(instance, tErr, log)
			return
		}
	}
	if watcher, ok := r.(WatchReconciler); ok {
		log.Info("Started watching resources...")
		if tErr := watcher.RefreshWatch(request); tErr != nil {
			log.Error(err, "Failed watching resources...")
			result, err = r.ManageError(instance, tErr, log)
			return
		}
	}
	if statusReconciler, ok := r.(StatusReconciler); ok {
		log.Info("Started updating status of instance...")
		if tErr := statusReconciler.UpdateStatus(instance); tErr != nil {
			log.Error(err, "Failed updating status...")
			result, err = r.ManageError(instance, tErr, log)
			return
		}
	}

	result, err = r.ManageSuccess(instance, log)
	return
}
