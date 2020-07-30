package resources

import (
	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("resource")

type Resource struct {
	Object controllerutil.Object
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

func (irs InstanceResourceSet) GetResources() Resources {
	return irs.resources
}

// func (res Resource) ToLockedResource() (*lockedresource.LockedResource, error) {
// 	err := res.fillGroupVersionKind()
// 	if err != nil {
// 		return nil, err
// 	}
// 	mapObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(res.Object)
// 	if err != nil {
// 		return nil, err
// 	}
// 	unstructuredObj := unstructured.Unstructured{Object: mapObj}
// 	return &lockedresource.LockedResource{Unstructured: unstructuredObj, ExcludedPaths: []string{".metadata", ".status"}}, nil
// }

// func (res Resources) ToLockedResources() ([]lockedresource.LockedResource, error) {
// 	tranformer := func(resource Resource) (lockedresource.LockedResource, error) {
// 		result, err := resource.ToLockedResource()
// 		return *result, err
// 	}
// 	result, err := fun.MapWithError(tranformer, res)
// 	return result.([]lockedresource.LockedResource), err
// }

func (res Resources) ToSlice() []Resource {
	return []Resource(res)
}

func (res Resources) GetObjects() []controllerutil.Object {
	objects := make([]controllerutil.Object, 0)
	for _, item := range res {
		objects = append(objects, item.Object)
	}
	return objects
}

func (res Resources) GetK8SResources() []resource.KubernetesResource {
	objects := make([]resource.KubernetesResource, len(res))
	for index, item := range res {
		objects[index] = item.Object
	}
	return objects
}

// func (irs InstanceResourceSet) ToLockedResources() ([]lockedresource.LockedResource, error) {
// 	return irs.resources.ToLockedResources()
// }

func (irs *InstanceResourceSet) addResource(newResources ...Resource) {
	irs.resources = append(irs.resources, newResources...)
}

func (res *Resource) fillGroupVersionKind() error {
	scheme := runtime.NewScheme()
	err := v1.AddToScheme(scheme)
	err = appsv1.AddToScheme(scheme)
	if err != nil {
		return err
	}
	gvks, _, err := scheme.ObjectKinds(res.Object)
	log.Info("Finding ObjectKinds", "ObjectKinds", gvks)
	if err != nil {
		return err
	}
	res.Object.GetObjectKind().SetGroupVersionKind(gvks[0])
	log.Info("SetGroupVersionKind to Object", "GroupVersionKind", res.Object.GetObjectKind().GroupVersionKind())
	return nil
}
