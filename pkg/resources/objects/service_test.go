package objects

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Configmap tests", func() {
	var (
		xrootd    *xrootdv1alpha1.XrootdCluster
		component types.ComponentName
	)
	BeforeEach(func() {
		xrootd = &xrootdv1alpha1.XrootdCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "instance",
			},
		}
	})
	Describe("Worker component", func() {
		BeforeEach(func() {
			component = constant.XrootdWorker
		})
		It("generates xrootd service", func() {
			name := utils.GetObjectName(component, xrootd.Name)
			service := GenerateXrootdService(xrootd, name, utils.GetComponentLabels(component, xrootd.Name), component)
			Expect(service.Name).Should(Equal(string(name)))
			By("check it is ClusterIP")
			Expect(service.Spec.Type).Should(Equal(corev1.ServiceTypeClusterIP))
			Expect(service.Spec.ClusterIP).Should(Equal(corev1.ClusterIPNone))
			Expect(service.Spec.Ports).Should(HaveLen(2))
			Expect(service.Spec.Ports[0].Port).Should(Equal(int32(constant.CmsdPort)))
			Expect(service.Spec.Ports[1].Port).Should(Equal(int32(constant.XrootdPort)))
		})
	})
	Describe("Redirector component", func() {
		BeforeEach(func() {
			component = constant.XrootdRedirector
		})
		It("generates xrootd service", func() {
			name := utils.GetObjectName(component, xrootd.Name)
			service := GenerateXrootdService(xrootd, name, utils.GetComponentLabels(component, xrootd.Name), component)
			Expect(service.Name).Should(Equal(string(name)))
			By("check it is ClusterIP")
			Expect(service.Spec.Type).Should(Equal(corev1.ServiceTypeClusterIP))
			Expect(service.Spec.ClusterIP).Should(Equal(corev1.ClusterIPNone))
			Expect(service.Spec.Ports).Should(HaveLen(2))
			Expect(service.Spec.Ports[0].Port).Should(Equal(int32(constant.CmsdPort)))
			Expect(service.Spec.Ports[1].Port).Should(Equal(int32(constant.XrootdPort)))
		})
	})
})
