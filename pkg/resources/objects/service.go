package objects

import (
	"github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	types "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateXrootdService generates xrootd service
func GenerateXrootdService(
	xrootd *v1alpha1.XrootdCluster, objectName types.ObjectName,
	compLabels types.Labels, componentName types.ComponentName,
) corev1.Service {
	name := string(objectName)
	labels := compLabels

	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: xrootd.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Port:     int32(constant.CmsdPort),
					Protocol: corev1.ProtocolTCP,
					Name:     string(constant.Cmsd),
				},
				{
					Port:     int32(constant.XrootdPort),
					Protocol: corev1.ProtocolTCP,
					Name:     string(constant.Xrootd),
				},
			},
			Selector: labels,
		},
	}
}
