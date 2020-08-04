package v1alpha1

import (
	"time"

	v1 "k8s.io/api/core/v1"
)

// ClusterPhase represents one of the runtime phase of the cluster
type ClusterPhase string

// ClusterConditionType represents on one of the runtime condition of the cluster
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

	// Condition keeps track of all cluster conditions, if they exist.
	Conditions []ClusterCondition `json:"conditions,omitempty"`

	// CurrentXrootdProtocol tracks the currently-used xrootd protocol info
	CurrentXrootdProtocol XrootdProtocolStatus `json:"currentXrootdProtocol"`

	RedirectorStatus MemberStatus `json:"redirectorStatus"`
	WorkerStatus     MemberStatus `json:"workerStatus"`
}

// MemberStatus defines the observed status of Xrootd member (worker/redirector)
type MemberStatus struct {
	// Size is the current size of the cluster
	Size int       `json:"size"`
	Pods PodStatus `json:"pods"`
}

// PodStatus defines the status of each of the member Pods for the specific component of xrootd cluster
type PodStatus struct {
	// Ready are the xrootd members that are ready to serve requests
	// The member names are the same as the xrootd pod names
	Ready []string `json:"ready"`
	// Unready are the xrootd members not ready to serve requests
	Unready []string `json:"unready"`
}

// ClusterCondition represents one current condition of the xrootd cluster.
type ClusterCondition struct {
	// Type of cluster condition.
	Type ClusterConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

func NewMemberStatus(ready []string, unready []string) MemberStatus {
	size := len(ready) + len(unready)
	return MemberStatus{
		Size: size,
		Pods: PodStatus{
			Ready:   ready,
			Unready: unready,
		},
	}
}

func (cs *XrootdStatus) SetPhase(p ClusterPhase) {
	cs.Phase = p
}

func (cs *XrootdStatus) SetReadyCondition() {
	c := newClusterCondition(ClusterConditionAvailable, v1.ConditionTrue, "Cluster available", "")
	cs.setClusterCondition(*c)
}

func (cs *XrootdStatus) ClearCondition(t ClusterConditionType) {
	pos, _ := getClusterCondition(cs, t)
	if pos == -1 {
		return
	}
	cs.Conditions = append(cs.Conditions[:pos], cs.Conditions[pos+1:]...)
}

func (cs *XrootdStatus) setClusterCondition(c ClusterCondition) {
	pos, cp := getClusterCondition(cs, c.Type)
	if cp != nil &&
		cp.Status == c.Status && cp.Reason == c.Reason && cp.Message == c.Message {
		return
	}

	if cp != nil {
		cs.Conditions[pos] = c
	} else {
		cs.Conditions = append(cs.Conditions, c)
	}
}

func getClusterCondition(status *XrootdStatus, t ClusterConditionType) (int, *ClusterCondition) {
	for i, c := range status.Conditions {
		if t == c.Type {
			return i, &c
		}
	}
	return -1, nil
}

func newClusterCondition(condType ClusterConditionType, status v1.ConditionStatus, reason, message string) *ClusterCondition {
	now := time.Now().Format(time.RFC3339)
	return &ClusterCondition{
		Type:               condType,
		Status:             status,
		LastUpdateTime:     now,
		LastTransitionTime: now,
		Reason:             reason,
		Message:            message,
	}
}

// XrootdProtocolStatus defines the version info and image of running xrootd software
type XrootdProtocolStatus struct {
	// Version is the current xrootd version
	Version string `json:"version"`
	// Image is the currently used image for xrootd containers
	Image string `json:"image"`
}
