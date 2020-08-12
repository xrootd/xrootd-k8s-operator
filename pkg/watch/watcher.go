package watch

import (
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Watcher defines a logic to execute on reconciliation request
// It's similar to Reconciler loop but it's unmanaged
type Watcher interface {
	// Watch will be called with request data
	Watch(requests <-chan reconcile.Request) error
}

var log = logf.Log.WithName("Watcher")
