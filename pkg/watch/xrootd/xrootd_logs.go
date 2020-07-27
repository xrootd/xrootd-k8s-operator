package xrootd

import (
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/msoap/byline"
	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/k8sutil"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8stypes "k8s.io/apimachinery/pkg/types"
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

const waitMemberReadyDelay = 5 * time.Second

func (lw LogsWatcher) Watch(requests <-chan reconcile.Request) error {
	var reqLogger logr.Logger
	for request := range requests {
		reqLogger = log.WithValues("request", request, "component", lw.Component)

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
	if err := lw.reconciler.GetClient().Status().Update(context.TODO(), instance); err != nil {
		return errors.Wrap(err, "failed updating xrootd status")
	}
	return nil
}

func (lw LogsWatcher) getXrootdOwnedStatefulSet(request reconcile.Request) (*appsv1.StatefulSet, error) {
	ss := &appsv1.StatefulSet{}
	ssName := k8stypes.NamespacedName{
		Namespace: request.Namespace,
		Name:      string(utils.GetObjectName(lw.Component, request.Name)),
	}
	if err := lw.reconciler.GetClient().Get(context.TODO(), ssName, ss); err != nil {
		return nil, errors.Wrap(err, "failed to get the statefulset")
	}
	return ss, nil
}

func (lw LogsWatcher) obtainLogsOfAllPods(request reconcile.Request, instance *xrootdv1alpha1.Xrootd, resultChannel chan<- bool) {
	donePods := 0
	var totalPods int
	if lw.Component == constant.XrootdRedirector {
		totalPods = int(instance.Spec.Redirector.Replicas)
	} else if lw.Component == constant.XrootdWorker {
		totalPods = int(instance.Spec.Worker.Replicas)
	}
	for {
		if donePods == totalPods {
			break
		}
		time.Sleep(waitMemberReadyDelay)
		ss, err := lw.getXrootdOwnedStatefulSet(request)
		if err != nil {
			continue
		}
		readyPods := int(ss.Status.ReadyReplicas)
		for i := donePods; i < readyPods; i++ {

		}
	}
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

	var err error
	var reader io.ReadCloser
	for {
		reader, err = k8sutil.GetPodLogStream(*pod, opt, clientset)
		if err != nil {
			if strings.Contains(err.Error(), "ContainerCreating") {
				reqLogger.V(1).Info("Container not started yet, retrying...", "error", err)
			} else {
				reqLogger.Error(err, "unable to get pod stream", "options", opt)
				resultChannel <- false
				return
			}
		} else {
			break
		}
	}
	defer reader.Close()

	lineReader := byline.NewReader(reader)

	var regex *regexp.Regexp
	if lw.Component == constant.XrootdRedirector {
		regex = regexp.MustCompile(`Protocol: redirector..+ logged in.$`)
	} else if lw.Component == constant.XrootdWorker {
		regex = regexp.MustCompile(`Protocol: Logged into .+$`)
	}

	reqLogger.Info("Grepping and reading...", "regex", regex)
	buffer := make([]byte, 50)
	read, err := lineReader.GrepByRegexp(regex).Read(buffer)
	reqLogger.V(1).Info("Read to buffer", "length", read, "buffer", buffer)

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
	if err = lw.reconciler.GetClient().Status().Update(context.TODO(), pod); err != nil {
		reqLogger.Error(err, "failed updating pod status", "status", pod.Status)
		resultChannel <- false
	}

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
