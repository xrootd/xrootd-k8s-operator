package objects

import (
	"github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InstanceVolumeSet struct {
	volumes types.Volumes
	meta    metav1.ObjectMeta
}

func newInstanceVolumeSet(meta metav1.ObjectMeta) *InstanceVolumeSet {
	return &InstanceVolumeSet{
		volumes: types.Volumes(make(map[string]v1.Volume)),
		meta:    meta,
	}
}

func getXrootdVolumeMounts(component types.ComponentName) []v1.VolumeMount {
	volumeMounts := []v1.VolumeMount{}
	return volumeMounts
}
