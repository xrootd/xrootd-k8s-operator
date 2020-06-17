package objects

import (
	"path/filepath"

	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/shivanshs9/xrootd-operator/pkg/utils"
	types "github.com/shivanshs9/xrootd-operator/pkg/utils/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type XrootdTemplateData struct{}

func getConfigMapName(controllerName string, container types.ContainerName, suffix string) string {
	return utils.SuffixName(controllerName, string(container), suffix)
}

func scanDir(root string, tmplData interface{}) map[string]string {
	files := map[string]string{}
	return files
}

func GenerateContainerConfigMap(
	xrootd *v1alpha1.Xrootd, objectName types.ObjectName,
	compLabels types.Labels, container types.ContainerName,
	subpath string,
) v1.ConfigMap {
	name := getConfigMapName(xrootd.Name, container, subpath)
	labels := compLabels
	tmplData := XrootdTemplateData{}
	rootDir := filepath.Join("/", "configmap", string(container), subpath)
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
