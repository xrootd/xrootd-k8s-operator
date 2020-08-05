package v1alpha1

import (
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XrootdVersionSpec defines the desired state of XrootdVersion
type XrootdVersionSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

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
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XrootdVersion is the Schema for the xrootdversions API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xrootdversions,scope=Namespaced
type XrootdVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XrootdVersionSpec   `json:"spec,omitempty"`
	Status XrootdVersionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XrootdVersionList contains a list of XrootdVersion
type XrootdVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XrootdVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XrootdVersion{}, &XrootdVersionList{})
}
