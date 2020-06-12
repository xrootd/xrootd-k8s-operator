package types

import (
	v1 "k8s.io/api/core/v1"
)

type Volumes map[string]v1.Volume

type Containers = []v1.Container

func (v Volumes) ToSlice() []v1.Volume {
	var volumes []v1.Volume
	for _, vol := range v {
		volumes = append(volumes, vol)
	}
	return volumes
}
