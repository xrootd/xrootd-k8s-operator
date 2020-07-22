package v1alpha1

type ClusterPhase string
type ClusterConditionType string

const (
	ClusterPhaseNone     ClusterPhase = ""
	ClusterPhaseCreating              = "Creating"
	ClusterPhaseRunning               = "Running"
	ClusterPhaseFailed                = "Failed"

	ClusterConditionAvailable  ClusterConditionType = "Available"
	ClusterConditionRecovering                      = "Recovering"
	ClusterConditionScaling                         = "Scaling"
	ClusterConditionUpgrading                       = "Upgrading"
)

// XrootdStatus defines the observed state of Xrootd
// +k8s:openapi-gen=true
type XrootdStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Phase is the cluster running phase
	Phase  ClusterPhase `json:"phase"`
	Reason string       `json:"reason,omitempty"`

	// CurrentVersion is the current xrootd version
	CurrentVersion string `json:"currentVersion"`

	RedirectorStatus MemberStatus `json:"redirectorStatus"`
	WorkerStatus     MemberStatus `json:"workerStatus"`
}

// MemberStatus defines the observed status of Xrootd member (worker/redirector)
type MemberStatus struct {
	// Size is the current size of the cluster
	Size int `json:"size"`
	// Ready are the xrootd members that are ready to serve requests
	// The member names are the same as the xrootd pod names
	Ready []string `json:"ready,omitempty"`
	// Unready are the xrootd members not ready to serve requests
	Unready []string `json:"unready,omitempty"`
}
