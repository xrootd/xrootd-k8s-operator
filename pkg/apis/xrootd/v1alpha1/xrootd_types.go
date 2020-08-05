package v1alpha1

import (
	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/catalog/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XrootdSpec defines the desired state of Xrootd
type XrootdSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Version must be a Xrootd version to use in the cluster pods.
	// The requested version must be installed in the target namespace
	// using XrootdVersion CRD.
	Version    types.CatalogVersion `json:"version"`
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

func (xrootd *Xrootd) SetVersionInfo(version catalogv1alpha1.XrootdVersion) {
	xrootd.Status.CurrentXrootdProtocol = XrootdProtocolStatus{
		Version: string(version.Spec.Version),
		Image:   version.Spec.Image,
	}
}
