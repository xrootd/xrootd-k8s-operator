package types

import (
	"strconv"
)

// ComponentName denotes distinct functional components in the architecture
type ComponentName string

// ConfigName is a special type of component used for ConfigMaps and Secrets
type ConfigName = ComponentName

// ContainerName denotes lowest-level components, i.e. containers
type ContainerName string

// ObjectName denotes distinct k8s resources using ObjectMeta.Name
type ObjectName string

// VolumeName denotes individual volume names to be later mounted in the containers
type VolumeName string

// ContainerPort denotes port in the container
type ContainerPort int

func (port *ContainerPort) String() string {
	return strconv.Itoa(int(*port))
}

// CatalogVersion denotes the version string of CRD
type CatalogVersion string
