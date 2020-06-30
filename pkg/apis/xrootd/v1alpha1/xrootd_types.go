package v1alpha1

import (
	"github.com/redhat-cop/operator-utils/pkg/util/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XrootdSpec defines the desired state of Xrootd
type XrootdSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Worker     XrootdWorkerSpec     `json:"worker,omitempty"`
	Redirector XrootdRedirectorSpec `json:"redirector,omitempty"`
	Config     XrootdConfigSpec     `json:"config,omitempty"`
}

// XrootdStorageSpec defines the storage spec of Xrootd workers
type XrootdStorageSpec struct {
	// Class must be a storage class
	Class string `json:"class,omitempty"`
	// Capacity must be a storage capacity and should be a valid quantity (ex, 1Gi)
	Capacity string `json:"capacity,omitempty"`
}

// XrootdWorkerSpec defines the desired state of Xrootd workers
type XrootdWorkerSpec struct {
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Image must have a tag
	// +kubebuilder:validation:Pattern=".+:.+"
	Image   string            `json:"image,omitempty"`
	Storage XrootdStorageSpec `json:"storage,omitempty"`
}

// XrootdRedirectorSpec defines the desired state of Xrootd redirectors
type XrootdRedirectorSpec struct {
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
	// Image must have a tag
	// +kubebuilder:validation:Pattern=".+:.+"
	Image string `json:"image,omitempty"`
}

// XrootdConfigSpec defines the config spec used to generate xrootd.cf
type XrootdConfigSpec struct {
}

// XrootdStatus defines the observed state of Xrootd
// +k8s:openapi-gen=true
type XrootdStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	apis.EnforcingReconcileStatus `json:",inline"`
}

// GetEnforcingReconcileStatus provides the EnforcingReconcileStatus
// func (xrootd *Xrootd) GetEnforcingReconcileStatus() apis.EnforcingReconcileStatus {
// 	return xrootd.Status.EnforcingReconcileStatus
// }

// // SetEnforcingReconcileStatus sets the EnforcingReconcileStatus
// func (xrootd *Xrootd) SetEnforcingReconcileStatus(reconcileStatus apis.EnforcingReconcileStatus) {
// 	xrootd.Status.EnforcingReconcileStatus = reconcileStatus
// }

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Xrootd is the Schema for the xrootds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xrootds,scope=Namespaced
type Xrootd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XrootdSpec   `json:"spec,omitempty"`
	Status XrootdStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XrootdList contains a list of Xrootd
type XrootdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Xrootd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Xrootd{}, &XrootdList{})
}
