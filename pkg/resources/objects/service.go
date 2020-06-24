package objects

import (
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	types "github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateXrootdService(
	xrootd *v1alpha1.Xrootd, objectName types.ObjectName,
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
					Port:     int32(constant.XrootdPort),
					Protocol: corev1.ProtocolTCP,
					Name:     string(constant.Xrootd),
				},
				{
					Port:     int32(constant.CmsdPort),
					Protocol: corev1.ProtocolTCP,
					Name:     string(constant.Cmsd),
				},
			},
			Selector: labels,
		},
	}
}
