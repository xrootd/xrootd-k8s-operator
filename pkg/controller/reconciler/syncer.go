package reconciler

import "k8s.io/apimachinery/pkg/runtime"

// SyncReconciler syncs k8s resources for a requested CR instance
type SyncReconciler interface {
	SyncResources(instance resource) error
	GetOwnedResourceKinds(instance runtime.Object) []runtime.Object
}
