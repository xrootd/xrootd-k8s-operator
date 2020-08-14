package xrootd

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch/xrootd"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *XrootdClusterReconciler) RefreshWatch(request reconcile.Request) error {
	reqLogger := r.Log.WithValues("xrootdcluster", request.NamespacedName)
	reqLogger.Info("Watching Xrootd resources...")
	return r.WatchManager.RefreshWatch(request)
}

func (r *XrootdClusterReconciler) AddXrootdLogger() {
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdRedirector, r))
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdWorker, r))
}
