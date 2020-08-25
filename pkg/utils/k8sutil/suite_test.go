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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
)

const SuiteName = "K8S Utilities"

func TestK8SUtilities(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		SuiteName,
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	testFramework.Start(func(cfg *rest.Config) client.Client {
		k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
		framework.ExpectNoError(err)
		return k8sClient
	})

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	testFramework.TeardownCluster()
})
