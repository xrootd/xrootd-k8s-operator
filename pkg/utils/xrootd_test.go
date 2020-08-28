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

package utils

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/catalog/v1alpha1"
)

var _ = Describe("Xrootd Util Tests", func() {
	// sets up test framework
	testFramework.InitOnRunningSuite()

	Describe("catalog.XrootdVersion", func() {
		var (
			versionCr   *catalogv1alpha1.XrootdVersion
			versionName string
		)

		BeforeEach(func() {
			versionName = "test-version"
			versionCr = &catalogv1alpha1.XrootdVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      versionName,
					Namespace: testFramework.GetNamespace(),
				},
				Spec: catalogv1alpha1.XrootdVersionSpec{
					Version: types.CatalogVersion("test"),
					Image:   "image:test",
				},
			}
		})

		AfterEach(func() {
			_ = testFramework.Client.Delete(context.TODO(), versionCr)
		})

		JustBeforeEach(func() {
			framework.ExpectNoError(testFramework.Client.Create(context.TODO(), versionCr))
			time.Sleep(2 * time.Second)
		})

		Context("when no XrootdVersion CR found with given name", func() {
			BeforeEach(func() {
				versionName = "404-version"
			})
			It("returns error", func() {
				_, err := GetXrootdVersionInfo(testFramework.Client, testFramework.GetNamespace(), versionName)
				framework.ExpectError(err)
			})
		})

		Context("when XrootdVersion CR found with given name", func() {
			It("returns version instance", func() {
				result, err := GetXrootdVersionInfo(testFramework.Client, testFramework.GetNamespace(), versionName)
				framework.ExpectNoError(err)
				Expect(result.Spec.Version).Should(Equal(types.CatalogVersion("test")))
			})
		})
	})
})
