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
	"context"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Pod Logs Tests", func() {
	var (
		pod        *corev1.Pod
		logOptions *corev1.PodLogOptions
	)

	// sets up test framework
	testFramework.InitOnRunningSuite()

	BeforeEach(func() {
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      framework.RandomAlphabaticalString(10),
				Namespace: testFramework.GetNamespace(),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "sleeping",
						Image:   "busybox",
						Command: []string{"sleep 3600"},
					},
				},
			},
		}
		logOptions = &corev1.PodLogOptions{
			Container: "sleeping",
		}
	})

	AfterEach(func() {
		_ = testFramework.Client.Delete(context.TODO(), pod)
	})

	JustBeforeEach(func() {
		framework.ExpectNoError(testFramework.Client.Create(context.TODO(), pod))
	})

	// Follow option doesn't actually generate a never-ending stream of logs right now
	// have to implement a infinite loop to keep requesting K8S API for new stream
	PContext("with follow true", func() {
		BeforeEach(func() {
			logOptions.Follow = true
		})
		It("gets never-ending stream of logs", func() {
			stream, err := GetPodLogStream(*pod, logOptions, testFramework.ClientSet)
			framework.ExpectNoError(err)
			Consistently(func() error {
				buffer := make([]byte, 20)
				_, err := stream.Read(buffer)
				return err
			}, "1s").ShouldNot(HaveOccurred())
		})
	})

	Context("with follow false", func() {
		It("gets bounded stream of logs", func() {
			stream, err := GetPodLogStream(*pod, logOptions, testFramework.ClientSet)
			framework.ExpectNoError(err)
			// ioutil.ReadAll returns no error if EOF has reached.
			// So for every bounded stream, it'll actually return true
			Consistently(func() error {
				_, err := ioutil.ReadAll(stream)
				return err
			}, "1s").ShouldNot(HaveOccurred())
		})
	})

	Context("with wrong container", func() {
		BeforeEach(func() {
			logOptions.Container = "wrong-container"
		})
		It("gets nil stream and request error", func() {
			_, err := GetPodLogStream(*pod, logOptions, testFramework.ClientSet)
			framework.ExpectError(err)
		})
	})
})
