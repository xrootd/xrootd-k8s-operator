package reconciler

import "k8s.io/apimachinery/pkg/runtime"

type SyncReconciler interface {
	SyncResources(instance resource) error
	GetOwnedResourceKinds(instance runtime.Object) []runtime.Object
}
