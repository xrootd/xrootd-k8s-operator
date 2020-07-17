package objects

import (
	corev1 "k8s.io/api/core/v1"
)

func getExecProbe(command []string, timeout int32) *corev1.Probe {
	handler := corev1.Handler{
		Exec: &corev1.ExecAction{
			Command: command,
		},
	}
	return &corev1.Probe{
		Handler:        handler,
		TimeoutSeconds: timeout,
	}
}
