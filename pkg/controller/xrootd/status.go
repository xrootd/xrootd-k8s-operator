package xrootd

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileXrootd) UpdateStatus(instance controllerutil.Object) error {
	xrootd := instance.(*xrootdv1alpha1.Xrootd)

	// Worker nodes
	workerSize := int(xrootd.Spec.Worker.Replicas)
	unreadyPods := make([]string, workerSize)
	workerPrefName := utils.GetObjectName(constant.XrootdWorker, xrootd.Name)
	for i := 0; i < workerSize; i++ {
		unreadyPods[i] = utils.SuffixName(string(workerPrefName), strconv.Itoa(i))
	}
	xrootd.Status.WorkerStatus = xrootdv1alpha1.NewMemberStatus([]string{}, unreadyPods)

	// Redirector nodes
	redirectorSize := int(xrootd.Spec.Redirector.Replicas)
	unreadyPods = make([]string, redirectorSize)
	redirectorPrefName := utils.GetObjectName(constant.XrootdRedirector, xrootd.Name)
	for i := 0; i < redirectorSize; i++ {
		unreadyPods[i] = utils.SuffixName(string(redirectorPrefName), strconv.Itoa(i))
	}
	xrootd.Status.RedirectorStatus = xrootdv1alpha1.NewMemberStatus([]string{}, unreadyPods)

	xrootd.Status.Phase = xrootdv1alpha1.ClusterPhaseCreating

	if err := r.GetClient().Status().Update(context.TODO(), xrootd); err != nil {
		return errors.Wrap(err, "failed updating xrootd instance status")
	}
	return nil
}
