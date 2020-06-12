package resources

import (
	"github.com/redhat-cop/operator-utils/pkg/util/lockedresourcecontroller/lockedresource"
	"github.com/shivanshs9/ty/fun"
	"github.com/shivanshs9/xrootd-operator/pkg/apis/xrootd/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("resource")

type Resource struct {
	Object runtime.Object
}

type Resources []Resource

// InstanceResourceSet contains Resources for a given Xrootd instance
type InstanceResourceSet struct {
	resources Resources
	xrootd    *v1alpha1.Xrootd
}

func NewInstanceResourceSet(xrootd *v1alpha1.Xrootd) *InstanceResourceSet {
	return &InstanceResourceSet{
		resources: Resources(make([]Resource, 0)),
		xrootd:    xrootd,
	}
}

func (res Resource) ToLockedResource() (*lockedresource.LockedResource, error) {
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

func (irs InstanceResourceSet) ToLockedResources() ([]lockedresource.LockedResource, error) {
	return irs.resources.ToLockedResources()
}

func (irs *InstanceResourceSet) addResource(newResources ...Resource) {
	log.Info("Adding resources...", "resources", newResources)
	irs.resources = append(irs.resources, newResources...)
}
