package objects

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Container tests", func() {
	var (
		instance  *xrootdv1alpha1.XrootdCluster
		component types.ComponentName
	)
	BeforeEach(func() {
		instance = &xrootdv1alpha1.XrootdCluster{
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
		}
	})
	AfterEach(func() {
	})

	Context("with xrootd image set on status", func() {
		BeforeEach(func() {
			instance.Status.CurrentXrootdProtocol.Image = "xrootd:latest"
		})

		Describe("Worker component", func() {
			BeforeEach(func() {
				component = constant.XrootdWorker
			})
			It("generate 2 containers, 4 mounts and 3 volumes", func() {
				containers, volumes := getXrootdContainersAndVolume(instance, component)
				By("checking containers")
				Expect(containers).Should(HaveLen(2))
				cmsd := containers[0]
				xrootd := containers[1]
				Expect(cmsd.Name).Should(Equal(string(constant.Cmsd)))
				Expect(xrootd.Name).Should(Equal(string(constant.Xrootd)))
				Expect(xrootd.Ports).Should(HaveLen(1))
				Expect(xrootd.Ports[0].ContainerPort).Should(Equal(int32(constant.XrootdPort)))
				Expect(cmsd.Ports).Should(HaveLen(0))

				By("checking volumes")
				Expect(volumes).Should(HaveLen(3))
				Expect(volumes[0].ConfigMap).ToNot(BeNil())
				Expect(volumes[1].ConfigMap).ToNot(BeNil())
				Expect(volumes[2].EmptyDir).ToNot(BeNil())
				Expect(volumes[2].Name).Should(Equal(string(constant.XrootdSharedAdminPathVolumeName)))

				By("checking xrootd container has PVC mounted")
				Expect(xrootd.VolumeMounts).Should(HaveLen(4))
				Expect(xrootd.VolumeMounts[3].Name).Should(Equal(getDataPVName(instance.Name)))
				Expect(xrootd.VolumeMounts[3].MountPath).Should(Equal("/data"))
			})
		})

		Describe("Redirector component", func() {
			BeforeEach(func() {
				component = constant.XrootdRedirector
			})
			It("generates 2 containers, 3 mounts and 3 volumes", func() {
				containers, volumes := getXrootdContainersAndVolume(instance, component)
				By("checking containers")
				Expect(containers).Should(HaveLen(2))
				cmsd := containers[0]
				xrootd := containers[1]
				Expect(cmsd.Name).Should(Equal(string(constant.Cmsd)))
				Expect(xrootd.Name).Should(Equal(string(constant.Xrootd)))
				Expect(xrootd.Ports).Should(HaveLen(1))
				Expect(xrootd.Ports[0].ContainerPort).Should(Equal(int32(constant.XrootdPort)))
				Expect(cmsd.Ports).Should(HaveLen(1))
				Expect(cmsd.Ports[0].ContainerPort).Should(Equal(int32(constant.CmsdPort)))

				By("checking volumes")
				Expect(volumes).Should(HaveLen(3))
				Expect(volumes[0].ConfigMap).ToNot(BeNil())
				Expect(volumes[1].ConfigMap).ToNot(BeNil())
				Expect(volumes[2].EmptyDir).ToNot(BeNil())
				Expect(volumes[2].Name).Should(Equal(string(constant.XrootdSharedAdminPathVolumeName)))
			})
		})
	})

	Context("with missing xrootd image on status", func() {
		It("fails to create container", func() {
			Expect(func() {
				getXrootdContainersAndVolume(instance, component)
			}).To(Panic())
		})
	})
})
