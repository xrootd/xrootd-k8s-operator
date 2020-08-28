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
	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var resources1 = []resource.KubernetesResource{
	&appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "first-sts",
		},
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
	&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cm",
		},
		Data: map[string]string{
			"hello": "world",
		},
	},
}

var resources2 = []resource.KubernetesResource{
	&appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "second-sts",
		},
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
	&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cm",
		},
		Data: map[string]string{
			"hello": "andrew",
		},
	},
}

func printDelta(delta compare.ResourceDelta, t *testing.T) {
	t.Logf("created: %v, updated: %v, removed: %v", len(delta.Added), len(delta.Updated), len(delta.Removed))
	if len(delta.Added) > 0 {
		t.Log("created:")
		for _, item := range delta.Added {
			t.Logf("- %v", item)
		}
	}
	if len(delta.Updated) > 0 {
		t.Log("updated:")
		for _, item := range delta.Updated {
			t.Logf("- %v", item)
		}
	}
	if len(delta.Removed) > 0 {
		t.Log("removed:")
		for _, item := range delta.Removed {
			t.Logf("- %v", item)
		}
	}
}

func TestComparatorForDifferentResources(t *testing.T) {
	comparator := GetComparator()
	delta := comparator.Comparator.CompareArrays(resources1, resources2)
	printDelta(delta, t)
	if !delta.HasChanges() {
		t.Errorf("detected same: %v", delta)
	}
	if len(delta.Added) != 1 || delta.Added[0].GetName() != "second-sts" {
		t.Errorf("wrong added")
	}
	if len(delta.Updated) != 1 || delta.Updated[0].GetName() != "test-cm" {
		t.Errorf("wrong updated")
	}
	if len(delta.Removed) != 1 || delta.Removed[0].GetName() != "first-sts" {
		t.Errorf("wrong removed")
	}
}

func TestComparatorForSameResources(t *testing.T) {
	comparator := GetComparator()
	delta := comparator.Comparator.CompareArrays(resources1, resources1)
	printDelta(delta, t)
	if delta.HasChanges() {
		t.Errorf("wrongly detected changes: %v", delta)
	}
}
