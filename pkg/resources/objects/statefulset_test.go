package objects

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Statefulset tests", func() {
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
			Spec: xrootdv1alpha1.XrootdClusterSpec{
				Redirector: xrootdv1alpha1.XrootdRedirectorSpec{
					Replicas: 2,
				},
				Worker: xrootdv1alpha1.XrootdWorkerSpec{
					Replicas: 3,
				},
			},
			Status: xrootdv1alpha1.XrootdClusterStatus{
				CurrentXrootdProtocol: xrootdv1alpha1.XrootdProtocolStatus{
					Image: "xrootd:latest",
				},
			},
		}
	})
	AfterEach(func() {
	})

	Describe("Worker component", func() {
		BeforeEach(func() {
			component = constant.XrootdWorker
		})
		Context("without storage spec", func() {
			It("fails creation", func() {
				Expect(func() {
					GenerateXrootdStatefulSet(xrootd, utils.GetObjectName(component, xrootd.Name), utils.GetComponentLabels(component, xrootd.Name), component)
				}).To(Panic())
			})
		})

		Context("with storage spec", func() {
			BeforeEach(func() {
				xrootd.Spec.Worker.Storage = xrootdv1alpha1.XrootdStorageSpec{
					Class:    "ephemeral",
					Capacity: "2G",
				}
			})
			It("generates a statefulset", func() {
				name := utils.GetObjectName(component, xrootd.Name)
				sts := GenerateXrootdStatefulSet(xrootd, name, utils.GetComponentLabels(component, xrootd.Name), component)
				Expect(sts.Name).Should(Equal(string(name)))
				Expect(sts.Spec.Template.Spec.Containers).Should(HaveLen(2))
				Expect(*sts.Spec.Replicas).Should(Equal(xrootd.Spec.Worker.Replicas))
			})
		})
	})

	Describe("Redirector component", func() {
		BeforeEach(func() {
			component = constant.XrootdRedirector
		})
		It("generates a statefulset", func() {
			name := utils.GetObjectName(component, xrootd.Name)
			sts := GenerateXrootdStatefulSet(xrootd, name, utils.GetComponentLabels(component, xrootd.Name), component)
			Expect(sts.Name).Should(Equal(string(name)))
			Expect(sts.Spec.Template.Spec.Containers).Should(HaveLen(2))
			Expect(*sts.Spec.Replicas).Should(Equal(xrootd.Spec.Redirector.Replicas))
		})
	})
})
