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

package framework

import (
	"path/filepath"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Framework supports common operations used by e2e tests; it will keep a client & a namespace for you.
type Framework struct {
	BaseName     string
	clientConfig *rest.Config
	Client       client.Client
	testEnv      *envtest.Environment

	Namespace *corev1.Namespace
	RootPath  string
}

// NewDefaultFramework creates a new Test Framework with CRDs
func NewDefaultFramework(baseName string, rootPath string) *Framework {
	logf.SetLogger(zap.LoggerTo(ginkgo.GinkgoWriter, true))

	rootPath, err := filepath.Abs(rootPath)
	ExpectNoError(err)

	testEnv := &envtest.Environment{
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths:     []string{filepath.Join(rootPath, "config", "crd", "bases")},
	}

	f := &Framework{
		BaseName: baseName,
		testEnv:  testEnv,
		RootPath: rootPath,
	}

	return f
}

// Start bootstraps the test env and sets the client
func (f *Framework) Start(getClient func(cfg *rest.Config) client.Client) {
	ginkgo.By("bootstrapping test environment")
	cfg, err := f.testEnv.Start()
	ExpectNoError(err)
	gomega.Expect(cfg).ToNot(gomega.BeNil())
	f.clientConfig = cfg
	f.Client = getClient(cfg)
	gomega.Expect(f.Client).ToNot(gomega.BeNil())
}

func (f *Framework) BeforeEach() {

}

// TeardownCluster stops the test environment
func (f *Framework) TeardownCluster() {
	ExpectNoError(f.testEnv.Stop())
}
