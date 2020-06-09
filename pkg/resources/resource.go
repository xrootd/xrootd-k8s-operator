package resources

import (
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	"github.com/shivanshs9/ty/fun"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource struct {
	Object runtime.Object
}

type Resources []Resource

func (res *Resource) ToLockedResource() (*lockedresource.LockedResource, error) {
	mapObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(res.Object)
	if err != nil {
		return nil, err
	}
	unstructuredObj := unstructured.Unstructured{Object: mapObj}
	return &lockedresource.LockedResource{Unstructured: unstructuredObj}, nil
}

func (resources Resources) ToLockedResources() ([]lockedresource.LockedResource, error) {
	tranformer := func(resource Resource) (lockedresource.LockedResource, error) {
		result, err := resource.ToLockedResource()
		return *result, err
	}
	result, err := fun.MapWithError(tranformer, resources)
	return result.([]lockedresource.LockedResource), err
}
