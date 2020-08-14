package objects

import (
	"fmt"
	"path/filepath"

	"github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
)

func getXrootdContainersAndVolume(xrootd *v1alpha1.XrootdCluster, component types.ComponentName) (types.Containers, types.Volumes) {
	volumeSet := newInstanceVolumeSet(xrootd.ObjectMeta)
	volumeSet.addEtcConfigVolume(constant.CfgXrootd)
	volumeSet.addRunConfigVolume(constant.CfgXrootd)
	// Shared filesystem is mounted for both xrootd and cmsd
	// Required for worker to communicate to cmsd using the named pipe located at 'adminpath'
	volumeSet.addEmptyDirVolume(constant.XrootdSharedAdminPathVolumeName, constant.XrootdSharedAdminPath)
	image := xrootd.Status.CurrentXrootdProtocol.Image
	if component == constant.XrootdWorker {
		volumeSet.addDataPVVolumeMount(filepath.Join("/", "data"))
	}
	volumeMounts := volumeSet.volumeMounts.ToSlice()

	probe := getExecProbe(
		[]string{
			"xrdfs",
			fmt.Sprintf("%s:%d", "127.0.0.1", constant.XrootdPort),
			"query",
			"config",
			"cms",
		},
		20)

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
			LivenessProbe: probe,
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
