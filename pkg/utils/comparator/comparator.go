/*


Copyright (C) 2020  The XRootD Collaboration

This library is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 2.1 of the License, or (at your option) any later version.

This library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public
License along with this library; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
USA
*/

package comparator

import (
	"reflect"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"k8s.io/apimachinery/pkg/api/equality"
)

const specField = "Spec"

var equals = equality.Semantic.DeepEqual
var zeroValue = reflect.Value{}

func deepEqual(existing resource.KubernetesResource, requested resource.KubernetesResource) bool {
	struct1 := reflect.ValueOf(existing).Elem()
	struct2 := reflect.ValueOf(requested).Elem()
	if spec1 := struct1.FieldByName(specField); spec1 != zeroValue {
		if spec2 := struct2.FieldByName(specField); spec2 != zeroValue {
			return equals(spec1.Addr().Interface(), spec2.Addr().Interface())
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
