/*


Copyright (C) 2020  The XRootD Collaboration

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
USA
*/

package controllers

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *XrootdClusterReconciler) UpdateStatus(instance controllerutil.Object) error {
	xrootd := instance.(*xrootdv1alpha1.XrootdCluster)

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
