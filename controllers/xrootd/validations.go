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
	"fmt"

	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsValid determines if a Xrootd instance is valid and initializes empty fields.
func (r *XrootdClusterReconciler) IsValid(instance controllerutil.Object) (result bool, err error) {
	result = true
	var invalidField string
	xrootd := instance.(*xrootdv1alpha1.XrootdCluster)
	// apply defaults
	r.applySpecDefaults(xrootd)

	// check capacity is valid quantity
	if _, tErr := resource.ParseQuantity(xrootd.Spec.Worker.Storage.Capacity); tErr != nil {
		result, err = false, errors.Wrapf(tErr, "Unable to parse storage capacity: '%v'", xrootd.Spec.Worker.Storage.Capacity)
		invalidField = ".Worker.Storage.Capacity"
	}

	if result {
		// check version is valid
		result, err = r.checkVersionIsValid(xrootd)
		if !result {
			invalidField = ".Version"
		}
	}

	if result {
		xrootd.Status.SetSpecValidCondition(true, "All fields valid", "'IsValid' check passed")
	} else {
		// if invalid spec, set the culprit field as the reason and error as message
		xrootd.Status.SetSpecValidCondition(false, invalidField, err.Error())
	}
	return
}

// applySpecDefaults sets the default values to fields
func (r *XrootdClusterReconciler) applySpecDefaults(xrootd *xrootdv1alpha1.XrootdCluster) {
	if xrootd.Spec.Redirector.Replicas == 0 {
		xrootd.Spec.Redirector.Replicas = 1
	}
	if xrootd.Spec.Worker.Replicas == 0 {
		xrootd.Spec.Worker.Replicas = 1
	}
	if len(xrootd.Spec.Worker.Storage.Class) == 0 {
		xrootd.Spec.Worker.Storage.Class = "standard"
	}
}

// checkVersionIsValid checks whether the version provided in spec is valid and respective version info is found
func (r *XrootdClusterReconciler) checkVersionIsValid(xrootd *xrootdv1alpha1.XrootdCluster) (result bool, err error) {
	result = true
	if len(xrootd.Spec.Version) == 0 {
		result, err = false, fmt.Errorf("provide xrootd version in instance")
	} else if versionInfo, tErr := utils.GetXrootdVersionInfo(r.GetClient(), xrootd.GetNamespace(), xrootd.Spec.Version); tErr != nil {
		result, err = false, errors.Wrapf(tErr, "unable to find requested version - %s", xrootd.Spec.Version)
	} else if image := versionInfo.Spec.Image; len(image) == 0 {
		result, err = false, fmt.Errorf("invalid image, '%s', provided for the given version, '%s'", image, xrootd.Spec.Version)
	} else {
		xrootd.SetVersionInfo(*versionInfo)
	}
	return
}
