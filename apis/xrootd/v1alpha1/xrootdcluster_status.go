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
	"time"

	v1 "k8s.io/api/core/v1"
)

// ClusterPhase represents one of the runtime phase of the cluster
type ClusterPhase string

// ClusterConditionType represents on one of the runtime condition of the cluster
type ClusterConditionType string

/*
These are valid Cluster Phases.
"ClusterPhaseInvalid" means the cluster spec is invalid.
"ClusterPhaseCreating" means the the cluster is being created.
"ClusterPhaseRunning" means the cluster is running in healthy state.
"ClusterPhaseFailed" means the cluster is failing.
*/
const (
	ClusterPhaseNone     ClusterPhase = ""
	ClusterPhaseCreating ClusterPhase = "Creating"
	ClusterPhaseRunning  ClusterPhase = "Running"
	ClusterPhaseFailed   ClusterPhase = "Failed"
)

/*
These are valid Cluster Condition types.
"ClusterConditionValid" means the cluster spec is valid.
"ClusterConditionAvailable" means the cluster is available to communicate.
"ClusterConditionRecovering" means the cluster is in recovering condition
"ClusterConditionScaling" means the cluster is scaling up/down.
"ClusterConditionUpgrading" means the cluster is undergoing a version upgrade.
*/
const (
	ClusterConditionValid      ClusterConditionType = "Valid"
	ClusterConditionAvailable  ClusterConditionType = "Available"
	ClusterConditionRecovering ClusterConditionType = "Recovering"
	ClusterConditionScaling    ClusterConditionType = "Scaling"
	ClusterConditionUpgrading  ClusterConditionType = "Upgrading"
)

// XrootdClusterStatus defines the observed state of XrootdCluster
// +k8s:openapi-gen=true
type XrootdClusterStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Phase is the current phase of the cluster
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Phase",xDescriptors="urn:alm:descriptor:io.kubernetes.phase"
	Phase ClusterPhase `json:"phase"`
	// Reason explains the current phase of the cluster.
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Phase Details",xDescriptors="urn:alm:descriptor:io.kubernetes.phase:reason"
	Reason string `json:"reason,omitempty"`

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

// NewMemberStatus creates a new xrootd member status with given ready and unready pods
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

// SetPhase sets the current phase of the cluster
func (cs *XrootdClusterStatus) SetPhase(p ClusterPhase) {
	cs.Phase = p
}

// SetReadyCondition sets the ClusterAvailable condition type to true
func (cs *XrootdClusterStatus) SetReadyCondition() {
	c := newClusterCondition(ClusterConditionAvailable, v1.ConditionTrue, "Cluster available", "")
	cs.setClusterCondition(*c)
}

// SetSpecValidCondition sets the ClusterValid condition type to given value
func (cs *XrootdClusterStatus) SetSpecValidCondition(isValid bool, reason string, msg string) {
	actualStatus := v1.ConditionUnknown
	if isValid {
		actualStatus = v1.ConditionTrue
	} else {
		actualStatus = v1.ConditionFalse
	}
	c := newClusterCondition(ClusterConditionValid, actualStatus, reason, msg)
	cs.setClusterCondition(*c)
}

// ClearCondition clears the given condition type
func (cs *XrootdClusterStatus) ClearCondition(t ClusterConditionType) {
	pos, _ := cs.GetClusterCondition(t)
	if pos == -1 {
		return
	}
	cs.Conditions = append(cs.Conditions[:pos], cs.Conditions[pos+1:]...)
}

func (cs *XrootdClusterStatus) setClusterCondition(c ClusterCondition) {
	pos, cp := cs.GetClusterCondition(c.Type)
	if cp != nil {
		cs.Conditions[pos] = c
	} else {
		cs.Conditions = append(cs.Conditions, c)
	}
}

// GetClusterCondition returns position and condition pointer from the .Conditions array
func (cs *XrootdClusterStatus) GetClusterCondition(t ClusterConditionType) (int, *ClusterCondition) {
	for i, c := range cs.Conditions {
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
	// Version is the current xrootd version used in the cluster
	Version string `json:"version"`
	// Image is the currently used image for xrootd containers
	Image string `json:"image"`
}
