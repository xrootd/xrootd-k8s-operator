package constant

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	XrootdPodConnection corev1.PodConditionType = corev1.PodConditionType("Connected")
)
