package comparator

import (
	"reflect"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"k8s.io/apimachinery/pkg/api/equality"
)

const specField = "Spec"

var equals = equality.Semantic.DeepEqual

func deepEqual(existing resource.KubernetesResource, requested resource.KubernetesResource) bool {
	struct1 := reflect.ValueOf(existing).Elem().Type()
	struct2 := reflect.ValueOf(requested).Elem().Type()
	if spec1, found1 := struct1.FieldByName(specField); found1 {
		if spec2, found2 := struct2.FieldByName(specField); found2 {
			return equals(spec1, spec2)
		}
	}
	return equals(existing, requested)
}

// GetComparator returns a MapComparator to compare k8s resources.
// It is useful when syncing resources to decide whether to
// create new, update or delete existing resources.
// This implementation uses equality.Semantic.DeepEqual on spec field of resources.
func GetComparator() *compare.MapComparator {
	comparator := compare.DefaultComparator()
	comparator.SetDefaultComparator(deepEqual)
	return &compare.MapComparator{Comparator: comparator}
}
