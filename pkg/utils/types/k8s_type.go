package types

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// Labels is synonymous to labels.Set
type Labels = labels.Set

// Volumes represents set of corev1.Volume
type Volumes []v1.Volume

// VolumeMounts represents set of corev1.VolumeMount
type VolumeMounts []v1.VolumeMount

// Containers represents set of corev1.Container
type Containers = []v1.Container

// ToSlice returns array of underlying volumes
func (v *Volumes) ToSlice() []v1.Volume {
	return []v1.Volume(*v)
}

// Add returns a new set of volumes by adding volumes to the existing set
func (v *Volumes) Add(volumes ...v1.Volume) Volumes {
	return append(v.ToSlice(), volumes...)
}

// ToSlice returns array of underlying volume mounts
func (vm *VolumeMounts) ToSlice() []v1.VolumeMount {
	return []v1.VolumeMount(*vm)
}

// Add returns a new set of volume mounts by adding volume mounts to the existing set
func (vm *VolumeMounts) Add(volumeMounts ...v1.VolumeMount) VolumeMounts {
	return append(vm.ToSlice(), volumeMounts...)
}
