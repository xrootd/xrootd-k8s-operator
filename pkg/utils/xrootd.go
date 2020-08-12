package utils

import (
	"context"

	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/catalog/v1alpha1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetXrootdVersionInfo returns the XrootdVersion CR object with Name as `versionName`
func GetXrootdVersionInfo(runtimeClient client.Client, namespace string, versionName string) (result *catalogv1alpha1.XrootdVersion, err error) {
	result = &catalogv1alpha1.XrootdVersion{}
	key := k8stypes.NamespacedName{
		Namespace: namespace,
		Name:      versionName,
	}
	err = runtimeClient.Get(context.TODO(), key, result)
	return
}
