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
	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/catalog/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XrootdClusterSpec defines the desired state of XrootdCluster
type XrootdClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Version must be name of XrootdVersion CR instance which defines the xrootd protcol image to use in the cluster pods.
	// The requested XrootdVersion instance must be installed in the target namespace using XrootdVersion CRD.
	Version    string               `json:"version"`
	Worker     XrootdWorkerSpec     `json:"worker,omitempty"`
	Redirector XrootdRedirectorSpec `json:"redirector,omitempty"`
	Config     XrootdConfigSpec     `json:"config,omitempty"`
}

// XrootdStorageSpec defines the storage spec of Xrootd workers
type XrootdStorageSpec struct {
	// Class must be a storage class
	// +kubebuilder:default=standard
	Class string `json:"class,omitempty"`
	// Capacity must be a storage capacity and should be a valid quantity (ex, 1Gi)
	Capacity string `json:"capacity,omitempty"`
}

// XrootdWorkerSpec defines the desired state of Xrootd workers
type XrootdWorkerSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32             `json:"replicas,omitempty"`
	Storage  XrootdStorageSpec `json:"storage,omitempty"`
}

// XrootdRedirectorSpec defines the desired state of Xrootd redirectors
type XrootdRedirectorSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`
}

// XrootdConfigSpec defines the config spec used to generate xrootd.cf
type XrootdConfigSpec struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// XrootdCluster is the Schema for the xrootdclusters API
type XrootdCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XrootdClusterSpec   `json:"spec,omitempty"`
	Status XrootdClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// XrootdClusterList contains a list of XrootdCluster
type XrootdClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XrootdCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XrootdCluster{}, &XrootdClusterList{})
}

// SetVersionInfo update the current version info of xrootd protocol
func (xrootd *XrootdCluster) SetVersionInfo(version catalogv1alpha1.XrootdVersion) {
	xrootd.Status.CurrentXrootdProtocol = XrootdProtocolStatus{
		Version: string(version.Spec.Version),
		Image:   version.Spec.Image,
	}
}
