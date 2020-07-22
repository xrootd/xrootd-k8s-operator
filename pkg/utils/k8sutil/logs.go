package k8sutil

import (
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodLogStream gives a IO stream of logs of pod with respect
// to provided options.
// Source - https://stackoverflow.com/a/53870271/2674983
func GetPodLogStream(pod corev1.Pod, opts *corev1.PodLogOptions, clientset *kubernetes.Clientset) (io.ReadCloser, error) {
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return nil, err
	}
	return podLogs, nil
}
