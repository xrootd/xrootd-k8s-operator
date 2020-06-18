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

func (ivs *InstanceVolumeSet) addVolumes(vols ...v1.Volume) {
	ivs.volumes = ivs.volumes.Add(vols...)
}

func (ivs *InstanceVolumeSet) addVolumeMounts(vols ...v1.VolumeMount) {
	ivs.volumeMounts = ivs.volumeMounts.Add(vols...)
}

func (ivs *InstanceVolumeSet) addConfigVolume(config types.ConfigName, subPath string, mountPath string, mode int32) {
	configmap := getConfigMapName(utils.GetObjectName(config, ivs.meta.Name), subPath)
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
	ivs.addVolumes(volume)
	ivs.addVolumeMounts(volumeMount)
}

func (ivs *InstanceVolumeSet) addEtcConfigVolume(config types.ConfigName) {
	ivs.addConfigVolume(config, "etc", "/config-etc", 0)
}

func (ivs *InstanceVolumeSet) addRunConfigVolume(config types.ConfigName) {
	ivs.addConfigVolume(config, "run", "/config-run", 0555)
}
