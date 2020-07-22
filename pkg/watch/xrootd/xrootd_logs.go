package xrootd

import (
	"github.com/msoap/byline"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/k8sutil"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type LogsWatcher struct {
	Component  types.ComponentName
	reconciler reconciler.Reconciler
}

var _ watch.Watcher = LogsWatcher{}

func (lw LogsWatcher) Watch(requests <-chan reconcile.Request) error {
	for request := range requests {

	}
	return nil
}

func (lw LogsWatcher) processXrootdLogs(request reconcile.Request) error {
	clientset, err := kubernetes.NewForConfig(lw.reconciler.GetConfig())
	if err != nil {
		return err
	}
	pods, err := lw.getXrootdOwnedPods(lw.Component, request)
	if err != nil {
		return err
	}
	opt := &corev1.PodLogOptions{
		Follow:    true,
		Container: string(constant.Cmsd),
	}
	for _, pod := range pods {
		reader, err := k8sutil.GetPodLogStream(pod, opt, clientset)
		if err != nil {
			return err
		}
		lineReader := byline.NewReader(reader)

		defer reader.Close()
	}
	return nil
}

func (lw LogsWatcher) getXrootdOwnedPods(component types.ComponentName, request reconcile.Request) ([]corev1.Pod, error) {

}

func NewLogsWatcher(component types.ComponentName, reconciler reconciler.Reconciler) watch.Watcher {
	return watch.NewGroupedRequestWatcher(
		LogsWatcher{
			Component:  component,
			reconciler: reconciler,
		},
	)
}
