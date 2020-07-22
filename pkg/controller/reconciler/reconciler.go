package reconciler

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type resource = controllerutil.Object

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

type BaseReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	config   *rest.Config
	recorder record.EventRecorder
}

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

func (br BaseReconciler) GetClient() client.Client {
	return br.client
}

func (br BaseReconciler) GetRecorder() record.EventRecorder {
	return br.recorder
}

func (br BaseReconciler) GetScheme() *runtime.Scheme {
	return br.scheme
}

func (br BaseReconciler) GetConfig() *rest.Config {
	return br.config
}

// IsBeingDeleted returns whether this object has been requested to be deleted
func IsBeingDeleted(obj resource) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

func (br *BaseReconciler) GetResourceInstance(request reconcile.Request, instance resource) error {
	if err := br.GetClient().Get(context.TODO(), request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
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
func Reconcile(r Reconciler, request reconcile.Request, instance resource, log logr.Logger) (reconcile.Result, error) {
	if err := r.GetResourceInstance(request, instance); err != nil {
		return r.ManageError(instance, err, log)
	}
	if ok, err := r.IsValid(instance); !ok {
		// If CR isn't valid, update the status with error and return
		return r.ManageError(instance, err, log)
	}

	if IsBeingDeleted(instance) {
		log.Info("Deleting instance...", "instance", instance)
		// TODO: Write Cleanup Logic
		return reconcile.Result{}, nil
	}

	if syncer, ok := r.(SyncReconciler); ok {
		log.Info("Starting syncing resources...")
		if err := syncer.SyncResources(instance); err != nil {
			log.Error(err, "Failing syncing resources...")
			return r.ManageError(instance, err, log)
		}
	}
	if watcher, ok := r.(WatchReconciler); ok {
		log.Info("Starting watching resources...")
		if err := watcher.RefreshWatch(request); err != nil {
			log.Error(err, "Failing watching resources...")
			return r.ManageError(instance, err, log)
		}
	}

	return r.ManageSuccess(instance, log)
}
