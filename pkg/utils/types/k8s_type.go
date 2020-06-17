package types

import (
	v1 "k8s.io/api/core/v1"
)

type Volumes []v1.Volume

type VolumeMounts []v1.VolumeMount

type Containers = []v1.Container

func (v *Volumes) ToSlice() []v1.Volume {
	return []v1.Volume(*v)
}

func (v *Volumes) Add(volumes ...v1.Volume) Volumes {
	return append(v.ToSlice(), volumes...)
}

func (vm *VolumeMounts) ToSlice() []v1.VolumeMount {
	return []v1.VolumeMount(*vm)
}

func (vm *VolumeMounts) Add(volumeMounts ...v1.VolumeMount) VolumeMounts {
	return append(vm.ToSlice(), volumeMounts...)
}
