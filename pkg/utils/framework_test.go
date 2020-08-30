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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/catalog/v1alpha1"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
)

var testFramework = framework.NewDefaultFramework("utils", framework.Options{
	UseManager: true,
})

const SuiteName = "K8S Utilities"

func TestUtilities(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		SuiteName,
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	var err error

	err = catalogv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)
	err = xrootdv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)

	testFramework.Start()

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	testFramework.TeardownCluster()
})
