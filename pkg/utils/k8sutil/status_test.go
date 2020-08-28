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

package k8sutil

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Pod Status Tests", func() {
	var (
		pod      *corev1.Pod
		condType corev1.PodConditionType
		podKey   types.NamespacedName
	)

	// sets up test framework
	testFramework.InitOnRunningSuite()

	BeforeEach(func() {
		condType = corev1.PodConditionType(framework.RandomAlphabaticalString(5))
		podKey = types.NamespacedName{
			Namespace: testFramework.GetNamespace(),
			Name:      fmt.Sprintf("%s-%s", "test-pod", framework.RandomAlphabaticalString(5)),
		}
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podKey.Name,
				Namespace: podKey.Namespace,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "sleeping",
						Image:   "busybox",
						Command: []string{"sleep 3600"},
					},
				},
			},
		}
	})

	AfterEach(func() {
		_ = testFramework.Client.Delete(context.TODO(), pod)
	})

	JustBeforeEach(func() {
		framework.ExpectNoError(testFramework.Client.Create(context.TODO(), pod))
	})

	initPodWithGivenConditionType := func(condType corev1.PodConditionType, condValue corev1.ConditionStatus) {
		now := metav1.NewTime(time.Now())
		pod.Status.Conditions = []corev1.PodCondition{
			{
				Type:               condType,
				Status:             condValue,
				Reason:             "??",
				LastProbeTime:      now,
				LastTransitionTime: now,
			},
		}
	}

	initPodWithArbitaryCondition := func() {
		initPodWithGivenConditionType(corev1.PodConditionType("arbitrary"), corev1.ConditionUnknown)
	}

	initPodWithSpecificCondition := func() {
		initPodWithGivenConditionType(condType, corev1.ConditionUnknown)
	}

	testUpdateConditionWithBool := func(condValue *bool) int {
		origLength := len(pod.Status.Conditions)
		By(fmt.Sprintf("update condition %s to %v", condType, condValue))
		framework.ExpectNoError(UpdatePodConditionWithBool(pod, condType, condValue, "test", testFramework.Client))
		return origLength
	}

	testAddNewCondition := func(condValue corev1.ConditionStatus, origConditionLength int) {
		By("fetching pod")
		fetched, err := testFramework.ClientSet.CoreV1().Pods(podKey.Namespace).Get(context.TODO(), podKey.Name, metav1.GetOptions{})
		framework.ExpectNoError(err)

		By("checking new condition added")
		conditions := fetched.Status.Conditions
		Expect(len(conditions)).Should(Equal(origConditionLength + 1))
		Expect(conditions[origConditionLength].Type).Should(Equal(condType))
		Expect(conditions[origConditionLength].Status).Should(Equal(condValue))
	}

	testUpdateExistingCondition := func(condValue corev1.ConditionStatus, origConditionLength int) {
		By("fetching pod")
		fetched, err := testFramework.ClientSet.CoreV1().Pods(podKey.Namespace).Get(context.TODO(), podKey.Name, metav1.GetOptions{})
		framework.ExpectNoError(err)

		By("checking existing condition updated")
		conditions := fetched.Status.Conditions
		Expect(len(conditions)).Should(Equal(origConditionLength))
		idx := origConditionLength - 1
		Expect(conditions[idx].Type).Should(Equal(condType))
		Expect(conditions[idx].Status).Should(Equal(condValue))
	}

	Describe("Update Pod Condition with bool", func() {
		Context("when the new condition status is nil", func() {
			Context("when the given pod doesn't have any condition", func() {
				It("would add a new condition", func() {
					testAddNewCondition(corev1.ConditionUnknown, testUpdateConditionWithBool(nil))
				})
			})

			Context("when the given pod have some other conditions", func() {
				BeforeEach(initPodWithArbitaryCondition)
				It("would add a new condition", func() {
					testAddNewCondition(corev1.ConditionUnknown, testUpdateConditionWithBool(nil))
				})
			})

			Context("when the given pod already has that condition", func() {
				BeforeEach(initPodWithSpecificCondition)
				It("would update the existing condition", func() {
					testUpdateExistingCondition(corev1.ConditionUnknown, testUpdateConditionWithBool(nil))
				})
			})
		})

		Context("when the new condition status is false", func() {
			Context("when the given pod doesn't have any condition", func() {
				It("would add a new condition", func() {
					condValue := false
					testAddNewCondition(corev1.ConditionFalse, testUpdateConditionWithBool(&condValue))
				})
			})

			Context("when the given pod have some other conditions", func() {
				BeforeEach(initPodWithArbitaryCondition)
				It("would add a new condition", func() {
					condValue := false
					testAddNewCondition(corev1.ConditionFalse, testUpdateConditionWithBool(&condValue))
				})
			})

			Context("when the given pod already has that condition", func() {
				BeforeEach(initPodWithSpecificCondition)
				It("would update the existing condition", func() {
					condValue := false
					testUpdateExistingCondition(corev1.ConditionFalse, testUpdateConditionWithBool(&condValue))
				})
			})
		})

		Context("when the new condition status is true", func() {
			Context("when the given pod doesn't have any condition", func() {
				It("would add a new condition", func() {
					condValue := true
					testAddNewCondition(corev1.ConditionTrue, testUpdateConditionWithBool(&condValue))
				})
			})

			Context("when the given pod have some other conditions", func() {
				BeforeEach(initPodWithArbitaryCondition)
				It("would add a new condition", func() {
					condValue := true
					testAddNewCondition(corev1.ConditionTrue, testUpdateConditionWithBool(&condValue))
				})
			})

			Context("when the given pod already has that condition", func() {
				BeforeEach(initPodWithSpecificCondition)
				It("would update the existing condition", func() {
					condValue := true
					testUpdateExistingCondition(corev1.ConditionTrue, testUpdateConditionWithBool(&condValue))
				})
			})
		})

		Context("when the given pod is nil", func() {
			It("panicks", func() {
				defer func() {
					if err, ok := recover().(error); ok {
						framework.ExpectError(err)
					}
				}()
				By(fmt.Sprintf("update condition %s to %v", condType, corev1.ConditionFalse))
				framework.ExpectNoError(UpdatePodCondition(nil, condType, corev1.ConditionFalse, "test", testFramework.Client))
			})
		})
	})
})
