package xrootd

import (
	"context"
	"regexp"

	"github.com/msoap/byline"
	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/k8sutil"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type LogsWatcher struct {
	Component  types.ComponentName
	reconciler reconciler.Reconciler
}

var _ watch.Watcher = LogsWatcher{}

var log = logf.Log.WithName("XrootdLogsWatcher")

func (lw LogsWatcher) Watch(requests <-chan reconcile.Request) error {
	for request := range requests {
		reqLogger := log.WithValues("request", request, "component", lw.Component)

		instance := &xrootdv1alpha1.Xrootd{}
		if err := lw.reconciler.GetResourceInstance(request, instance); err != nil {
			return err
		}
		if reconciler.IsBeingDeleted(instance) {
			// Skip processing if requested instance is being deleted
			reqLogger.Info("Xrootd instance is being deleted...", "request", request)
			continue
		}
		if err := lw.monitorXrootdStatus(request, instance); err != nil {
			reqLogger.Error(err, "Failed to monitor xrootd cluster...")
		}
	}
	return nil
}

func (lw LogsWatcher) monitorXrootdStatus(request reconcile.Request, instance *xrootdv1alpha1.Xrootd) error {
	reqLogger := log.WithValues("request", request, "component", lw.Component)
	reqLogger.Info("Started monitoring xrootd cluster...")

	clientset, err := kubernetes.NewForConfig(lw.reconciler.GetConfig())
	if err != nil {
		return errors.Wrap(err, "unable to get kubernetes clientset")
	}
	pods, err := lw.getXrootdOwnedPods(request)
	if err != nil {
		return errors.Wrap(err, "unable to get pods owned by the xrootd instance")
	}
	reqLogger.Info("Fetched pods...", "pods", len(pods.Items))
	opt := &corev1.PodLogOptions{
		Follow:    true,
		Container: string(constant.Cmsd),
	}
	resultChannel := make(chan bool)
	for _, pod := range pods.Items {
		go lw.processXrootdPodLogs(&pod, opt, clientset, resultChannel)
	}

	currPod := 0
	readyPods := make([]string, len(pods.Items))
	unreadyPods := make([]string, len(pods.Items))
	for isPodReady := range resultChannel {
		if isPodReady {
			readyPods = append(readyPods, pods.Items[currPod].Name)
		} else {
			unreadyPods = append(unreadyPods, pods.Items[currPod].Name)
		}
		currPod++
	}
	status := utils.NewMemberStatus(readyPods, unreadyPods)
	if lw.Component == constant.XrootdWorker {
		instance.Status.WorkerStatus = status
	} else if lw.Component == constant.XrootdRedirector {
		instance.Status.RedirectorStatus = status
	}
	return nil
}

func (lw LogsWatcher) getXrootdOwnedPods(request reconcile.Request) (*corev1.PodList, error) {
	pods := &corev1.PodList{}
	selector := labels.NewSelector()
	for key, value := range utils.GetComponentLabels(lw.Component, request.Name) {
		req, err := labels.NewRequirement(key, selection.Equals, []string{value})
		if err != nil {
			return nil, err
		}
		selector = selector.Add(*req)
	}
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     request.Namespace,
	}
	if err := lw.reconciler.GetClient().List(context.TODO(), pods, &opts); err != nil {
		return nil, errors.Wrap(err, "failed listing pods")
	}
	return pods, nil
}

func (lw LogsWatcher) processXrootdPodLogs(pod *corev1.Pod, opt *corev1.PodLogOptions, clientset *kubernetes.Clientset, resultChannel chan<- bool) {
	reqLogger := log.WithValues("pod", pod.Name, "component", lw.Component)

	reader, err := k8sutil.GetPodLogStream(*pod, opt, clientset)
	if err != nil {
		reqLogger.Error(err, "unable to get pod stream", "options", opt)
		resultChannel <- false
		return
	}
	defer reader.Close()

	lineReader := byline.NewReader(reader)

	var regex *regexp.Regexp
	if lw.Component == constant.XrootdRedirector {
		regex, err = regexp.Compile(`Protocol: redirector..+ logged in.$`)
	} else if lw.Component == constant.XrootdWorker {
		regex, err = regexp.Compile(`Protocol: Logged into .+$`)
	}
	if err != nil {
		log.Error(err, "regex compile error", "component", lw.Component)
		resultChannel <- false
		return
	}

	log.Info("Grepping and reading...", "regex", regex)
	buffer := make([]byte, 50)
	read, err := lineReader.GrepByRegexp(regex).Read(buffer)
	log.Info("Read to buffer", "length", read, "buffer", buffer)

	result := read > 0

	status := corev1.ConditionFalse
	if result {
		status = corev1.ConditionTrue
	}
	pod.Status.Conditions = append(pod.Status.Conditions, corev1.PodCondition{
		Type:   constant.XrootdPodConnection,
		Status: status,
		Reason: "Cmsd logs confirmed logged-in status",
	})

	resultChannel <- result
}

func NewLogsWatcher(component types.ComponentName, reconciler reconciler.Reconciler) watch.Watcher {
	return watch.NewGroupedRequestWatcher(
		LogsWatcher{
			Component:  component,
			reconciler: reconciler,
		},
	)
}
