package objects

import (
	"path/filepath"

	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
)

func getXrootdContainersAndVolume(xrootd *v1alpha1.Xrootd, component types.ComponentName) (types.Containers, types.Volumes) {
	spec := xrootd.Spec
	volumeSet := newInstanceVolumeSet(xrootd.ObjectMeta)
	volumeSet.addEtcConfigVolume(constant.CfgXrootd)
	volumeSet.addRunConfigVolume(constant.CfgXrootd)
	// Shared filesystem is mounted for both xrootd and cmsd
	// Required for worker to communicate to cmsd using the named pipe located at 'adminpath'
	volumeSet.addEmptyDirVolume(constant.XrootdSharedAdminPathVolumeName, constant.XrootdSharedAdminPath)
	var image string
	if component == constant.XrootdRedirector {
		image = spec.Redirector.Image
	} else {
		image = spec.Worker.Image
	}
	if component == constant.XrootdWorker {
		volumeSet.addDataPVVolumeMount(filepath.Join("/", "data"))
	}
	volumeMounts := volumeSet.volumeMounts.ToSlice()

	containers := []v1.Container{
		{
			Name:            string(constant.Cmsd),
			Image:           image,
			ImagePullPolicy: v1.PullIfNotPresent,
			VolumeMounts:    volumeMounts,
			Command:         constant.ContainerCommand,
			Args:            []string{"-s", "cmsd"},
			WorkingDir:      "/tmp",
		},
		{
			Name:            string(constant.Xrootd),
			Image:           image,
			ImagePullPolicy: v1.PullIfNotPresent,
			VolumeMounts:    volumeMounts,
			Command:         constant.ContainerCommand,
			Args:            []string{"-s", "xrootd"},
			WorkingDir:      "/tmp",
			Ports: []v1.ContainerPort{
				{
					Name:          string(constant.Xrootd),
					ContainerPort: int32(constant.XrootdPort),
					Protocol:      v1.ProtocolTCP,
				},
			},
		},
	}

	if component == constant.XrootdRedirector {
		containers[0].Ports = []v1.ContainerPort{
			{
				Name:          string(constant.Cmsd),
				ContainerPort: int32(constant.CmsdPort),
				Protocol:      v1.ProtocolTCP,
			},
		}
	}

	return containers, volumeSet.volumes
}
