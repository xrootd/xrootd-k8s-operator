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
	"context"
	"fmt"
	"path/filepath"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Framework supports common operations used by e2e tests; it will keep a client & a namespace for you.
type Framework struct {
	BaseName     string
	UniqueName   string
	clientConfig *rest.Config
	Client       client.Client
	ClientSet    clientset.Interface
	testEnv      *envtest.Environment

	SkipNamespaceCreation bool
	Namespace             *corev1.Namespace
	RootPath              string
}

// NewDefaultFramework creates a new Test Framework with CRDs
func NewDefaultFramework(baseName string, rootPath string) *Framework {
	logf.SetLogger(zap.LoggerTo(ginkgo.GinkgoWriter, true))

	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		panic(err)
	}

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
	Logf("bootstrapping test environment")
	cfg, err := f.testEnv.Start()
	ExpectNoError(err)
	gomega.Expect(cfg).ToNot(gomega.BeNil())
	f.clientConfig = cfg

	f.ClientSet, err = clientset.NewForConfig(cfg)
	ExpectNoError(err)
	gomega.Expect(f.ClientSet).ToNot(gomega.BeNil())

	f.Client = getClient(cfg)
	gomega.Expect(f.Client).ToNot(gomega.BeNil())
}

// InitOnRunningSuite sets up ginkgo's BeforeEach & AfterEach.
// It must be called within running ginkgo suite (like Describe, Context etc)
func (f *Framework) InitOnRunningSuite() {
	ginkgo.BeforeEach(f.beforeEach)
	ginkgo.AfterEach(f.afterEach)
}

// beforeEach sets up a random namespace if allowed
func (f *Framework) beforeEach() {
	if !f.SkipNamespaceCreation {
		ginkgo.By(fmt.Sprintf("Building a namespace api object, basename %s", f.BaseName))
		namespace, err := f.CreateNamespace(f.BaseName, map[string]string{
			"e2e-framework": f.BaseName,
		})
		ExpectNoError(err)

		f.Namespace = namespace
		f.UniqueName = namespace.GetName()
	} else {
		f.UniqueName = fmt.Sprintf("%s-%s", f.BaseName, RandomAlphabaticalString(8))
	}
}

// afterEach collect reports and cleans up after each test
func (f *Framework) afterEach() {

}

// CreateNamespace creates a new namespace with randomized name from the given baseName and labels
func (f Framework) CreateNamespace(baseName string, labels map[string]string) (*corev1.Namespace, error) {
	labels["e2e-run"] = string(RunID)
	name := fmt.Sprintf("%s-%s", baseName, RandomAlphabaticalString(8))
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "",
			Labels:    labels,
		},
	}
	var got *corev1.Namespace
	var err error
	maxAttempts := 3
	for attempt := 0; attempt < maxAttempts; attempt++ {
		got, err = f.ClientSet.CoreV1().Namespaces().Create(context.TODO(), namespaceObj, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				// regenerate on conflict
				Logf("Namespace name %q was already taken, generate a new name and retry", namespaceObj.Name)
				namespaceObj.Name = fmt.Sprintf("%v-%v", baseName, RandomAlphabaticalString(8))
			} else {
				Logf("Unexpected error while creating namespace: %v", err)
			}
		} else {
			break
		}
	}
	return got, err
}

// GetNamespace returns the ephemeral namespace for the test
func (f Framework) GetNamespace() string {
	if f.Namespace != nil {
		return f.Namespace.Name
	}
	return "default"
}

// TeardownCluster stops the test environment
func (f *Framework) TeardownCluster() {
	ExpectNoError(f.testEnv.Stop())
}
