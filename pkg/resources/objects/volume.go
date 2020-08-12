package objects

import (
	"github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func getDataPVName(name string) string {
	return utils.SuffixName(name, "data")
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

func (ivs *InstanceVolumeSet) addEmptyDirVolume(volumeName types.VolumeName, path string) {
	volume := v1.Volume{
		Name: string(volumeName),
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	volumeMount := v1.VolumeMount{
		Name:      string(volumeName),
		MountPath: path,
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

func (ivs *InstanceVolumeSet) addDataPVVolumeMount(mountPath string) {
	volumeMount := v1.VolumeMount{
		Name:      getDataPVName(ivs.meta.Name),
		MountPath: mountPath,
	}
	ivs.addVolumeMounts(volumeMount)
}

func getDataPVClaim(xrootd *v1alpha1.Xrootd) v1.PersistentVolumeClaim {
	defer func() {
		if err := recover(); err != nil {
			rLog.WithName("volume.DataPVClaim").Error(err.(error), "failed parsing storage capacity", "xrootd", xrootd)
		}
	}()
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: getDataPVName(xrootd.Name),
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			StorageClassName: &xrootd.Spec.Worker.Storage.Class,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": resource.MustParse(xrootd.Spec.Worker.Storage.Capacity),
				},
			},
		},
	}
}
