package objects

import (
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
)

func getXrootdContainersAndVolume(xrootd *v1alpha1.Xrootd, component types.ComponentName) (types.Containers, types.Volumes) {
	spec := xrootd.Spec
	volumeMounts := getXrootdVolumeMounts(component)
	var image string
	if component == constant.XrootdRedirector {
		image = spec.Redirector.Image
	} else {
		image = spec.Worker.Image
	}
	containers := []v1.Container{
		{
			Name:         string(constant.Cmsd),
			Image:        image,
			VolumeMounts: volumeMounts,
			Command:      []string{"echo hi"},
		},
		{
			Name:         string(constant.Xrootd),
			Image:        image,
			VolumeMounts: volumeMounts,
			Command:      []string{"echo hi"},
		},
	}

	volumeSet := newInstanceVolumeSet(xrootd.ObjectMeta)
	return containers, volumeSet.volumes
}
