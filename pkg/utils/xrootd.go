package utils

import (
	"context"

	"github.com/pkg/errors"
	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/catalog/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetFirstValidXrootdVersion(runtimeClient client.Client, namespace string, version types.CatalogVersion) (result *catalogv1alpha1.XrootdVersion, err error) {
	opt := &client.ListOptions{
		Namespace: namespace,
	}
	versionList := &catalogv1alpha1.XrootdVersionList{}
	if err = runtimeClient.List(context.TODO(), versionList, opt); err != nil {
		err = errors.Wrapf(err, "cannot find any xrootd version CRs in the %s namespace", namespace)
		return
	}
	for _, versionInfo := range versionList.Items {
		if versionInfo.Spec.Version == version {
			result = &versionInfo
			break
		}
	}
	return
}
