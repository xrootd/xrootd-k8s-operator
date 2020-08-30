package constant

import "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"

// ControllerName is the name of xrootd controller
const ControllerName = "xrootd-controller"

const (
	// XrootdRedirector is the name of xrootd redirector component
	XrootdRedirector types.ComponentName = "xrootd-redirector"
	// XrootdWorker is the name of xrootd worker component
	XrootdWorker types.ComponentName = "xrootd-worker"
)

const (
	// CfgXrootd is the name of xrootd config
	CfgXrootd types.ConfigName = "xrootd"
)

const (
	// Xrootd is the xrootd container name
	Xrootd types.ContainerName = "xrootd"
	// Cmsd is the cmsd container name
	Cmsd types.ContainerName = "cmsd"
)

// ControllerLabels is the default Labels for resources managed by this controller
var ControllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": ControllerName,
}

const (
	// XrootdPort is the xrootd container port
	XrootdPort types.ContainerPort = 1094
	// CmsdPort is the cmsd container port
	CmsdPort types.ContainerPort = 2131
)

// ContainerCommand is the run command for xrootd containers
var ContainerCommand = []string{"/config-run/start.sh"}

// XrootdSharedAdminPathVolumeName is the xrootd shared admin path volume name
const XrootdSharedAdminPathVolumeName types.VolumeName = "adminpath"

// XrootdSharedAdminPath is the mount path for the volume
const XrootdSharedAdminPath = "/tmp/xrd"

// EnvXrootdOpConfigmapPath is the environment key for path to configmaps folder
const EnvXrootdOpConfigmapPath = "XROOTD_OPERATOR_CONFIGMAPS_PATH"
