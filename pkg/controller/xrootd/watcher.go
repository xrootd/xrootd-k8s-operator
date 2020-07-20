package xrootd

import "sigs.k8s.io/controller-runtime/pkg/reconcile"

func (r *ReconcileXrootd) RefreshWatch(request reconcile.Request) error {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Watching Xrootd resources...")
	return r.WatchManager.RefreshWatch(request)
}
