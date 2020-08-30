package objects

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/template"
	types "github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type xrootdTemplateData struct {
	XrootdRedirectorDn       string
	XrootdRedirectorReplicas int
	XrootdPort               types.ContainerPort
	CmsdPort                 types.ContainerPort
	XrootdSharedPath         string
}

func getConfigMapName(objectName types.ObjectName, suffix string) string {
	return utils.SuffixName(string(objectName), suffix)
}

func scanDir(root string, tmplData interface{}) (map[string]string, error) {
	log := rLog.WithName("scanDir")
	files := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) (er error) {
		log.Info("Scanning file...", "path", path)
		if err == nil && !info.IsDir() {
			files[info.Name()], er = template.ApplyTemplate(path, tmplData)
		} else if err != nil {
			er = err
		}
		return
	})
	if err != nil {
		return nil, errors.Wrapf(err, "scanDir failed for '%s'", root)
	}
	return files, nil
}

// GenerateContainerConfigMap generated configmap for given xrootd container
func GenerateContainerConfigMap(
	xrootd *v1alpha1.XrootdCluster, objectName types.ObjectName,
	compLabels types.Labels, config types.ConfigName,
	subpath string,
) v1.ConfigMap {
	name := getConfigMapName(objectName, subpath)
	labels := compLabels
	var tmplData interface{}
	if config == constant.CfgXrootd {
		tmplData = xrootdTemplateData{
			XrootdRedirectorDn:       string(utils.GetObjectName(constant.XrootdRedirector, xrootd.Name)),
			XrootdRedirectorReplicas: int(xrootd.Spec.Redirector.Replicas),
			XrootdPort:               constant.XrootdPort,
			CmsdPort:                 constant.CmsdPort,
			XrootdSharedPath:         constant.XrootdSharedAdminPath,
		}
	}
	rootDir := os.Getenv(constant.EnvXrootdOpConfigmapPath)
	if len(rootDir) == 0 {
		rootDir = "configmaps"
	}
	rootPath, err := filepath.Abs(rootDir)
	if err != nil {
		panic(fmt.Errorf("error in getting absolute path, %v: %v", rootDir, err))
	}
	configDir := filepath.Join(rootPath, string(config), subpath)
	data, err := scanDir(configDir, tmplData)
	if err != nil {
		rLog.WithName("GenerateContainerConfigMap").Error(err, "Unable to apply template", "templateData", tmplData)
		panic(err)
	}
	return v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Labels:    labels,
			Namespace: xrootd.Namespace,
		},
		Data: data,
	}
}
