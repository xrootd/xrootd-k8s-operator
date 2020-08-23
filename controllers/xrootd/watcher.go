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

package xrootdcontroller

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch/xrootd"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// RefreshWatch implements WatchReconciler and runs RefreshWatch on the watch manager
func (r *XrootdClusterReconciler) RefreshWatch(request reconcile.Request) error {
	reqLogger := r.Log.WithValues("xrootdcluster", request.NamespacedName)
	reqLogger.Info("Watching Xrootd resources...")
	return r.WatchManager.RefreshWatch(request)
}

// AddXrootdLogger adds the Logs Watchers for Xrootd Redirector and Worker components
func (r *XrootdClusterReconciler) AddXrootdLogger() {
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdRedirector, r))
	r.AddWatchers(xrootd.NewLogsWatcher(constant.XrootdWorker, r))
}
