package constant

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// XrootdPodConnection indicates whether the cmsd container in the pod
	// is successfully connected to the xrootd protocol in the network.
	XrootdPodConnection corev1.PodConditionType = corev1.PodConditionType("Connected")
)
