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
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

var testFramework *framework.Framework
var k8sManager ctrl.Manager

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	testFramework = framework.NewDefaultFramework(ControllerName, filepath.Join("..", ".."))

	framework.ExpectNoError(os.Setenv(constant.EnvXrootdOpConfigmapPath, filepath.Join(testFramework.RootPath, "configmaps")))

	var err error
	err = catalogv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)
	err = xrootdv1alpha1.AddToScheme(scheme.Scheme)
	framework.ExpectNoError(err)

	// +kubebuilder:scaffold:scheme

	testFramework.Start(func(cfg *rest.Config) client.Client {
		k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
		framework.ExpectNoError(err)
		return k8sManager.GetClient()
	})

	// setup xrootd controller
	err = (&XrootdClusterReconciler{
		BaseReconciler: reconciler.NewBaseReconciler(
			k8sManager.GetClient(), k8sManager.GetScheme(),
			k8sManager.GetEventRecorderFor(ControllerName), k8sManager.GetConfig(),
		),
		WatchManager: reconciler.NewWatchManager(nil),
		Log:          ctrl.Log.WithName("controllers").WithName("XrootdCluster"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// start manager
	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(ctrl.SetupSignalHandler())).Should(Succeed())
	}()

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	testFramework.TeardownCluster()
})
