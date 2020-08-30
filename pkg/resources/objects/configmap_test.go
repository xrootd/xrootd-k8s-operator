package objects

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Configmap tests", func() {
	var (
		xrootd *xrootdv1alpha1.XrootdCluster
	)
	BeforeEach(func() {
		xrootd = &xrootdv1alpha1.XrootdCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "instance",
			},
			Spec: xrootdv1alpha1.XrootdClusterSpec{
				Redirector: xrootdv1alpha1.XrootdRedirectorSpec{
					Replicas: 2,
				},
			},
		}
	})
	AfterEach(func() {
		framework.ExpectNoError(os.Unsetenv(constant.EnvXrootdOpConfigmapPath))
	})

	Context("with configmap path env unset", func() {
		It("uses default and fails to create", func() {
			Expect(func() {
				GenerateContainerConfigMap(xrootd, utils.GetObjectName(constant.XrootdRedirector, xrootd.Name), utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name), constant.CfgXrootd, "dir")
			}).To(Panic())
		})
	})

	Context("with valid configmap path but invalid subpath", func() {
		It("fails to create", func() {
			os.Setenv(constant.EnvXrootdOpConfigmapPath, "../../../configmaps")
			Expect(func() {
				GenerateContainerConfigMap(xrootd, utils.GetObjectName(constant.XrootdRedirector, xrootd.Name), utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name), constant.CfgXrootd, "dir")
			}).To(Panic())
		})
	})

	Context("with valid configmap path and valid subpath", func() {
		It("generates a configmap", func() {
			os.Setenv(constant.EnvXrootdOpConfigmapPath, "../../../configmaps")
			name := utils.GetObjectName(constant.XrootdRedirector, xrootd.Name)
			cm := GenerateContainerConfigMap(xrootd, name, utils.GetComponentLabels(constant.XrootdRedirector, xrootd.Name), constant.CfgXrootd, "etc")
			Expect(cm.Name).Should(Equal(utils.SuffixName(string(name), "etc")))
			Expect(cm.Data).ShouldNot(BeEmpty())
		})
	})
})
