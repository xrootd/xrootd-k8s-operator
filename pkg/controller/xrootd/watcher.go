package xrootd

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch/xrootd"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileXrootd) RefreshWatch(request reconcile.Request) error {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Watching Xrootd resources...")
	return r.WatchManager.RefreshWatch(request)
}

func (r *ReconcileXrootd) AddXrootdLogger() {
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdRedirector, r))
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdWorker, r))
}
