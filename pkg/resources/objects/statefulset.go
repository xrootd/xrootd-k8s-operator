package objects

import (
	"github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateXrootdStatefulSet generated xrootd statefulset
func GenerateXrootdStatefulSet(
	xrootd *v1alpha1.XrootdCluster, objectName types.ObjectName,
	compLabels types.Labels, componentName types.ComponentName,
) appsv1.StatefulSet {
	labels := compLabels
	name := string(objectName)
	var replicas int32
	if componentName == constant.XrootdRedirector {
		replicas = xrootd.Spec.Redirector.Replicas
	} else {
		replicas = xrootd.Spec.Worker.Replicas
	}
	containers, volumes := getXrootdContainersAndVolume(xrootd, componentName)
	ss := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: xrootd.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: containers,
					Volumes:    volumes.ToSlice(),
				},
			},
		},
	}
	if componentName == constant.XrootdWorker {
		pvc, err := getDataPVClaim(xrootd)
		if err != nil {
			rLog.WithName("GenerateXrootdStatefulSet").Error(err, "could not create pvc for worker", "xrootd", xrootd)
			panic(err)
		}
		ss.Spec.VolumeClaimTemplates = []v1.PersistentVolumeClaim{*pvc}
	}
	return ss
}
