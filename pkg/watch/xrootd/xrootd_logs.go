package xrootd

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type LogsWatcher struct {
	watch.GroupedRequestWatcher
	Component types.ComponentName
}

var _ watch.Watcher = LogsWatcher{}

func NewLogsWatcher(component types.ComponentName) LogsWatcher {
	lw := &LogsWatcher{
		Component: component,
	}
	lw.GroupedRequestWatcher = watch.NewGroupedRequestWatcher(lw.processXrootdLogs)
	return *lw
}

func (lw LogsWatcher) processXrootdLogs(requestChannel <-chan reconcile.Request) {

}
