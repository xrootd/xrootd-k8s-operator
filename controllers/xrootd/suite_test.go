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
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	catalogv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/catalog/v1alpha1"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	framework.ExpectNoError(os.Setenv(constant.EnvXrootdOpConfigmapPath, filepath.Join(testFramework.GetRootPath(), "configmaps")))

	var err error
	err = catalogv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)
	err = xrootdv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)

	// +kubebuilder:scaffold:scheme

	testFramework.Start()

	// setup xrootd controller
	err = (&XrootdClusterReconciler{
		BaseReconciler: reconciler.NewBaseReconciler(
			testFramework.GetManager().GetClient(), testFramework.GetManager().GetScheme(),
			testFramework.GetManager().GetEventRecorderFor(ControllerName), testFramework.GetManager().GetConfig(),
		),
		WatchManager: reconciler.NewWatchManager(nil),
		Log:          ctrl.Log.WithName("controllers").WithName("XrootdCluster"),
	}).SetupWithManager(testFramework.GetManager())
	Expect(err).ToNot(HaveOccurred())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	testFramework.TeardownCluster()
})
