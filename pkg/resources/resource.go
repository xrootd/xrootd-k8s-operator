package resources

import (
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource struct {
	Object runtime.Object
}

func (res *Resource) ToLockedResource() (*lockedresource.LockedResource, error) {
	mapObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(res.Object)
	if err != nil {
		return nil, err
	}
	unstructuredObj := unstructured.Unstructured{Object: mapObj}
	return &lockedresource.LockedResource{Unstructured: unstructuredObj}, nil
}
