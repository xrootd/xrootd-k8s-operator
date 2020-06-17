package objects

import (
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstanceVolumeSet contains a set of v1.Volume for a given ObjectMeta
type InstanceVolumeSet struct {
	volumes      types.Volumes
	volumeMounts types.VolumeMounts
	meta         metav1.ObjectMeta
}

func getConfigVolumeName(configmapName string) string {
	return utils.SuffixName("config", configmapName)
}

func newInstanceVolumeSet(meta metav1.ObjectMeta) *InstanceVolumeSet {
	return &InstanceVolumeSet{
		volumes:      types.Volumes(make([]v1.Volume, 0)),
		volumeMounts: types.VolumeMounts(make([]v1.VolumeMount, 0)),
		meta:         meta,
	}
}

func (ivs *InstanceVolumeSet) addConfigVolume(container types.ContainerName, subPath string, mountPath string, mode int32) {
	configmap := getConfigMapName(ivs.meta.Name, container, subPath)
	volumeName := getConfigVolumeName(configmap)
	volume := v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configmap,
				},
			},
		},
	}
	volumeMount := v1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
	}
	if mode != 0 {
		volume.VolumeSource.ConfigMap.DefaultMode = &mode
	}
	ivs.volumes.Add(volume)
	ivs.volumeMounts.Add(volumeMount)
}

func (ivs *InstanceVolumeSet) addEtcConfigVolume(container types.ContainerName) {
	ivs.addConfigVolume(container, "etc", "/config-etc", 0)
}

func (ivs *InstanceVolumeSet) addRunConfigVolume(container types.ContainerName) {
	ivs.addConfigVolume(container, "run", "/config-run", 0555)
}
