package constant

import . "github.com/shivanshs9/xrootd-operator/pkg/utils/types"

const ControllerName = "xrootd-controller"

const (
	XrootdRedirector ComponentName = "xrootd-redirector"
	XrootdWorker     ComponentName = "xrootd-worker"
)

const (
	ConfigMap KindName = "cfg"
	Service   KindName = "svc"
)

const (
	Xrootd ContainerName = "xrootd"
)

var ControllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": ControllerName,
}
