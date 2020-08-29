package objects

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Volume", func() {
	var (
		instanceMeta metav1.ObjectMeta
	)
	BeforeEach(func() {
		instanceMeta = metav1.ObjectMeta{
			Namespace: "ns",
			Name:      "instance",
		}
	})

	Describe("InstanceVolumeSet", func() {
		var (
			ivs *InstanceVolumeSet
		)

		BeforeEach(func() {
			ivs = newInstanceVolumeSet(instanceMeta)
		})

		Context("when new volume and mount is added", func() {
			var (
				volume corev1.Volume
				mount  corev1.VolumeMount
			)

			BeforeEach(func() {
				volume = corev1.Volume{
					Name: "vol",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/apis",
						},
					},
				}
				mount = corev1.VolumeMount{
					Name:      "vol",
					MountPath: "/mnt",
				}
			})

			It("is added", func() {
				By("adding volume")
				origLength := len(ivs.volumes)
				ivs.addVolumes(volume)
				Expect(ivs.volumes).Should(HaveLen(origLength + 1))

				By("adding volume mount")
				origLength = len(ivs.volumeMounts)
				ivs.addVolumeMounts(mount)
				Expect(ivs.volumeMounts).Should(HaveLen(origLength + 1))
			})
		})

		Context("when adding empty dir", func() {
			It("both volume and mounts are added", func() {
				const volName string = "empty"
				const path string = "/tmp"
				ivs.addEmptyDirVolume(types.VolumeName(volName), path)
				By("check volume is added")
				newVol := ivs.volumes[0]
				Expect(newVol.Name).Should(Equal(volName))
				Expect(*newVol.EmptyDir).Should(Equal(corev1.EmptyDirVolumeSource{}))

				By("check mount is added")
				newMnt := ivs.volumeMounts[0]
				Expect(newMnt.Name).Should(Equal(volName))
				Expect(newMnt.MountPath).Should(Equal(path))
			})
		})

		Context("when adding config dirs", func() {
			It("both volume and mounts are added", func() {
				By("config-etc added")
				ivs.addEtcConfigVolume(constant.CfgXrootd)
				newVol := ivs.volumes[0]
				newMnt := ivs.volumeMounts[0]
				Expect(newVol.ConfigMap.Name).Should(Equal("instance-xrootd-etc"))
				Expect(newMnt.Name).Should(Equal("config-instance-xrootd-etc"))

				By("config-run added")
				ivs.addRunConfigVolume(constant.CfgXrootd)
				newVol = ivs.volumes[1]
				newMnt = ivs.volumeMounts[1]
				Expect(newVol.ConfigMap.Name).Should(Equal("instance-xrootd-run"))
				Expect(newMnt.Name).Should(Equal("config-instance-xrootd-run"))
				By("- check the file is executable")
				Expect(*newVol.ConfigMap.DefaultMode).Should(Equal(int32(0555)))
			})
		})

		Context("when adding PV", func() {
			It("only mounts is added", func() {
				ivs.addDataPVVolumeMount("/mnt")
				Expect(ivs.volumes).Should(HaveLen(0))
				Expect(ivs.volumeMounts).Should(HaveLen(1))
				newMnt := ivs.volumeMounts[0]
				Expect(newMnt.Name).Should(Equal("instance-data"))
				Expect(newMnt.MountPath).Should(Equal("/mnt"))
			})
		})
	})

	Describe("CR Data PV Claim", func() {
		var xrootd *xrootdv1alpha1.XrootdCluster

		BeforeEach(func() {
			xrootd = &xrootdv1alpha1.XrootdCluster{
				ObjectMeta: instanceMeta,
				Spec: xrootdv1alpha1.XrootdClusterSpec{
					Worker: xrootdv1alpha1.XrootdWorkerSpec{
						Storage: xrootdv1alpha1.XrootdStorageSpec{
							Class: "ephemeral",
						},
					},
				},
			}
		})
		Context("when capacity is invalid", func() {
			BeforeEach(func() {
				xrootd.Spec.Worker.Storage.Capacity = "sad"
			})
			It("it errors", func() {
				_, err := getDataPVClaim(xrootd)
				framework.ExpectError(err)
			})
		})
		Context("when capacity is valid", func() {
			BeforeEach(func() {
				xrootd.Spec.Worker.Storage.Capacity = "2Gi"
			})
			It("it returns a PVClaim", func() {
				pvc, err := getDataPVClaim(xrootd)
				framework.ExpectNoError(err)
				Expect(pvc.Name).Should(Equal("instance-data"))
				Expect(*pvc.Spec.StorageClassName).Should(Equal(xrootd.Spec.Worker.Storage.Class))
				Expect(pvc.Spec.Resources.Requests["storage"]).Should(Equal(resource.MustParse(xrootd.Spec.Worker.Storage.Capacity)))
			})
		})
	})
})
