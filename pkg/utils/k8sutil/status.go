package k8sutil

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdatePodCondition updates the given condition type of the given pod with the new status and reason
func UpdatePodCondition(pod *corev1.Pod, condType corev1.PodConditionType, condStatus corev1.ConditionStatus, reason string, kubeclient client.Client) error {
	now := metav1.NewTime(time.Now())
	found := false
	for i, condition := range pod.Status.Conditions {
		if condition.Type == condType {
			found = true
			condition.LastProbeTime = now
			if condition.Status != condStatus {
				condition.Status = condStatus
				condition.LastTransitionTime = now
			}
			pod.Status.Conditions[i] = condition
			break
		}
	}
	if !found {
		pod.Status.Conditions = append(pod.Status.Conditions, corev1.PodCondition{
			Type:               condType,
			Status:             condStatus,
			Reason:             reason,
			LastProbeTime:      now,
			LastTransitionTime: now,
		})
	}
	return kubeclient.Status().Update(context.TODO(), pod)
}

// UpdatePodConditionWithBool updates the given condition type of the given pod with the new status and reason.
// It accepts bool conditional status and internally calls UpdatePodCondition with correct status value
func UpdatePodConditionWithBool(pod *corev1.Pod, condType corev1.PodConditionType, condStatus *bool, reason string, kubeclient client.Client) error {
	actualStatus := corev1.ConditionUnknown
	if condStatus != nil {
		if *condStatus {
			actualStatus = corev1.ConditionTrue
		} else {
			actualStatus = corev1.ConditionFalse
		}
	}
	return UpdatePodCondition(pod, condType, actualStatus, reason, kubeclient)
}
