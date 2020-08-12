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

package v1alpha1

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XrootdVersionSpec defines the desired state of XrootdVersion
type XrootdVersionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Xrootd Version in the provided image
	Version types.CatalogVersion `json:"version"`
	// Whether this version is deprecated or not
	Deprecated bool `json:"deprecated,omitempty"`
	// Image must have a tag
	// +kubebuilder:validation:Pattern=".+:.+"
	Image string `json:"image"`
}

// XrootdVersionStatus defines the observed state of XrootdVersion
type XrootdVersionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// XrootdVersion is the Schema for the xrootdversions API
type XrootdVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XrootdVersionSpec   `json:"spec,omitempty"`
	Status XrootdVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// XrootdVersionList contains a list of XrootdVersion
type XrootdVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XrootdVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XrootdVersion{}, &XrootdVersionList{})
}
