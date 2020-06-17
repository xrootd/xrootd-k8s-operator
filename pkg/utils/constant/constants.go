package constant

import . "github.com/shivanshs9/xrootd-operator/pkg/utils/types"

const ControllerName = "xrootd-controller"

const (
	XrootdRedirector ComponentName = "xrootd-redirector"
	XrootdWorker     ComponentName = "xrootd-worker"
)

const (
	ConfigMap   KindName = "cfg"
	Service     KindName = "svc"
	StatefulSet KindName = "sts"
)

const (
	Xrootd ContainerName = "xrootd"
	Cmsd   ContainerName = "cmsd"
)

var ControllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": ControllerName,
}

const (
	XrootdPort ContainerPort = 1094
	CmsdPort   ContainerPort = 2131
)

var ContainerCommand = []string{"/config-run/start.sh"}
