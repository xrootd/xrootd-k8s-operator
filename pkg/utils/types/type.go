package types

import (
	"strconv"

	"k8s.io/apimachinery/pkg/labels"
)

type Labels = labels.Set

type ComponentName string

type ConfigName = ComponentName

type ContainerName string

type KindName string

type ObjectName string

type VolumeName string

type ContainerPort int

func (port *ContainerPort) String() string {
	return strconv.Itoa(int(*port))
}

type CatalogVersion string
