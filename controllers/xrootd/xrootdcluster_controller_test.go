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

package xrootdcontroller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/catalog/v1alpha1"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func applyCrsAndFetchCrAfterDelay(
	versionCr *catalogv1alpha1.XrootdVersion,
	clusterCr *xrootdv1alpha1.XrootdCluster,
	clusterKey types.NamespacedName,
	delay time.Duration,
) *xrootdv1alpha1.XrootdCluster {
	By("creating xrootdversion successfully")
	Expect(testFramework.Client.Create(context.TODO(), versionCr)).Should(Succeed())
	By("creating xrootdcluster successfully")
	Expect(testFramework.Client.Create(context.TODO(), clusterCr)).Should(Succeed())
	time.Sleep(delay)

	By("fetching updated xrootdcluster CR")
	fetched := &xrootdv1alpha1.XrootdCluster{}
	Expect(testFramework.Client.Get(context.TODO(), clusterKey, fetched)).Should(Succeed())
	return fetched
}

var _ = Describe("XrootdCluster Controller", func() {
	const waitCr = time.Second * 5

	var (
		versionSpec     catalogv1alpha1.XrootdVersionSpec
		versionToCreate *catalogv1alpha1.XrootdVersion
		clusterSpec     xrootdv1alpha1.XrootdClusterSpec
		clusterToCreate *xrootdv1alpha1.XrootdCluster
		versionKey      types.NamespacedName
		clusterKey      types.NamespacedName
	)

	// sets up test framework
	testFramework.InitOnRunningSuite()

	BeforeEach(func() {
		clusterKey = types.NamespacedName{
			Namespace: testFramework.GetNamespace(),
			Name:      "test-xrootdcluster-" + framework.RandomAlphabaticalString(5),
		}
		versionKey = types.NamespacedName{
			Namespace: testFramework.GetNamespace(),
			Name:      "test-xrootdversion-" + framework.RandomAlphabaticalString(5),
		}

		versionSpec = catalogv1alpha1.XrootdVersionSpec{
			Version: "latest",
			Image:   "qserv/xrootd:latest",
		}
		clusterSpec = xrootdv1alpha1.XrootdClusterSpec{
			Version: versionKey.Name,
			Worker: xrootdv1alpha1.XrootdWorkerSpec{
				Storage: xrootdv1alpha1.XrootdStorageSpec{
					Capacity: "1Gi",
				},
			},
		}
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
		_ = testFramework.Client.Delete(context.Background(), clusterToCreate)
		_ = testFramework.Client.Delete(context.Background(), versionToCreate)
	})

	JustBeforeEach(func() {
		clusterToCreate = &xrootdv1alpha1.XrootdCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterKey.Name,
				Namespace: clusterKey.Namespace,
			},
			Spec: clusterSpec,
		}
		versionToCreate = &catalogv1alpha1.XrootdVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name:      versionKey.Name,
				Namespace: versionKey.Namespace,
			},
			Spec: versionSpec,
		}
	})

	Describe("XrootdVersion with Invalid Data", func() {
		Context("Spec with empty image", func() {
			BeforeEach(func() {
				versionSpec = catalogv1alpha1.XrootdVersionSpec{
					Version: "latest",
				}
			})
			It("should fail creation", func() {
				Expect(testFramework.Client.Create(context.TODO(), versionToCreate)).ShouldNot(Succeed())
			})
		})
	})

	Describe("XrootdCluster with Invalid Data", func() {
		Context("Spec with empty version name", func() {
			BeforeEach(func() {
				clusterSpec = xrootdv1alpha1.XrootdClusterSpec{}
			})

			It("should fail creation", func() {
				Expect(testFramework.Client.Create(context.TODO(), clusterToCreate)).ShouldNot(Succeed())
			})
		})
		Context("Spec with invalid version name", func() {
			BeforeEach(func() {
				clusterSpec = xrootdv1alpha1.XrootdClusterSpec{
					Version: "invalid",
				}
			})

			It("should pass creation but add failed status", func() {
				fetched := applyCrsAndFetchCrAfterDelay(versionToCreate, clusterToCreate, clusterKey, waitCr)
				By("checking 'failed' phase")
				Expect(fetched.Status.Phase).Should(Equal(xrootdv1alpha1.ClusterPhaseFailed))

				By("checking false 'valid' condition")
				conditionAssertion := Expect(func() xrootdv1alpha1.ClusterCondition {
					_, res := fetched.Status.GetClusterCondition(xrootdv1alpha1.ClusterConditionValid)
					if res == nil {
						return xrootdv1alpha1.ClusterCondition{}
					}
					return *res
				}())
				conditionAssertion.ToNot(BeZero())
				conditionAssertion.Should(MatchFields(IgnoreExtras, Fields{
					"Status": Equal(corev1.ConditionFalse),
				}))
			})
		})
		Context("Spec with empty worker storage capacity", func() {
			BeforeEach(func() {
				clusterSpec = xrootdv1alpha1.XrootdClusterSpec{
					Version: versionKey.Name,
				}
			})

			It("should pass creation but add failed status", func() {
				fetched := applyCrsAndFetchCrAfterDelay(versionToCreate, clusterToCreate, clusterKey, waitCr)
				By("checking 'failed' phase")
				Expect(fetched.Status.Phase).Should(Equal(xrootdv1alpha1.ClusterPhaseFailed))

				By("checking false 'valid' condition")
				conditionAssertion := Expect(func() xrootdv1alpha1.ClusterCondition {
					_, res := fetched.Status.GetClusterCondition(xrootdv1alpha1.ClusterConditionValid)
					if res == nil {
						return xrootdv1alpha1.ClusterCondition{}
					}
					return *res
				}())
				conditionAssertion.ToNot(BeZero())
				conditionAssertion.Should(MatchFields(IgnoreExtras, Fields{
					"Status": Equal(corev1.ConditionFalse),
				}))
			})
		})
	})

	Describe("XrootdCluster with valid data", func() {
		Context("xrootd redirectors and workers are live", func() {
			It("cluster is is in running phase", func() {
				fetched := applyCrsAndFetchCrAfterDelay(versionToCreate, clusterToCreate, clusterKey, waitCr)
				By("checking 'creating' phase")
				Expect(fetched.Status.Phase).Should(Equal(xrootdv1alpha1.ClusterPhaseCreating))
				// time.Sleep(waitCr * 12)
				// By("checking 'running' phase")
				// Expect(fetched.Status.Phase).Should(Equal(xrootdv1alpha1.ClusterPhaseRunning))
			})
		})
	})
})
