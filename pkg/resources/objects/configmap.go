package objects

import (
	"os"
	"path/filepath"

	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/constant"
	"github.com/shivanshs9/xrootd-operator/pkg/utils/template"
	types "github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type xrootdTemplateData struct {
	XrootdRedirectorDn       string
	XrootdRedirectorReplicas int
	XrootdPort               types.ContainerPort
	CmsdPort                 types.ContainerPort
}

func getConfigMapName(objectName types.ObjectName, suffix string) string {
	return utils.SuffixName(string(objectName), suffix)
}

func scanDir(root string, tmplData interface{}) map[string]string {
	log := rLog.WithName("scanDir")
	files := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) (er error) {
		log.Info("Scanning file...", "path", path)
		if err == nil && !info.IsDir() {
			files[info.Name()], er = template.ApplyTemplate(path, tmplData)
		}
		return er
	})
	if err != nil {
		log.Error(err, "Unable to apply template for", "root", root, "templateData", tmplData)
	}
	return files
}

func GenerateContainerConfigMap(
	xrootd *v1alpha1.Xrootd, objectName types.ObjectName,
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
		}
	}
	rootDir := filepath.Join("/", "configmaps", string(config), subpath)
	data := scanDir(rootDir, tmplData)
	return v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Labels:    labels,
			Namespace: xrootd.Namespace,
		},
		Data: data,
	}
}
