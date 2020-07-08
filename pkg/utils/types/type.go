package types

import "strconv"

type Labels = map[string]string

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
