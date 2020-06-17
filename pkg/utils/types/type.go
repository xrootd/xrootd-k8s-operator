package types

import "strconv"

type Labels = map[string]string

type ComponentName string

type ContainerName string

type KindName string

type ObjectName string

type ContainerPort int

func (port *ContainerPort) String() string {
	return strconv.Itoa(int(*port))
}
