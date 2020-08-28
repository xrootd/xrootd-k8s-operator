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
	"testing"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var resources1 = []resource.KubernetesResource{
	&appsv1.StatefulSet{
		Spec: appsv1.StatefulSetSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "test",
							Name:  "dummy",
						},
					},
				},
			},
		},
	},
}

var resources2 = []resource.KubernetesResource{
	&appsv1.StatefulSet{
		Spec: appsv1.StatefulSetSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "test",
							Name:  "wow",
						},
					},
				},
			},
		},
	},
}

func TestComparatorForDifferentResources(t *testing.T) {
	comparator := GetComparator()
	delta := comparator.Comparator.CompareArrays(resources1, resources2)
	if !delta.HasChanges() {
		t.Errorf("detected same: %v", delta)
	}
}

func TestComparatorForSameResources(t *testing.T) {
	comparator := GetComparator()
	delta := comparator.Comparator.CompareArrays(resources1, resources1)
	if delta.HasChanges() {
		t.Errorf("wrongly detected changes: %v", delta)
	}
}
