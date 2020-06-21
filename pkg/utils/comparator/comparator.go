package comparator

import (
	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"k8s.io/apimachinery/pkg/api/equality"
)

func deepEqual(existing resource.KubernetesResource, requested resource.KubernetesResource) bool {
	return equality.Semantic.DeepEqual(existing, requested)
}

func GetComparator() compare.ResourceComparator {
	comparator := compare.DefaultComparator()
	comparator.SetDefaultComparator(deepEqual)
	return comparator
}
