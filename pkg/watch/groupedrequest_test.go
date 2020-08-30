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

package watch

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xrootd/xrootd-k8s-operator/tests/integration/framework"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("GroupedRequest watcher", func() {
	const waitPerRequest = 1 * time.Second

	var (
		grw      *GroupedRequestWatcher
		watcher  TestWatcher
		requests chan reconcile.Request
	)

	BeforeEach(func() {
		watcher = newTestWatcher(waitPerRequest)
		grw = NewGroupedRequestWatcher(watcher).(*GroupedRequestWatcher)
		requests = make(chan reconcile.Request, 3)
	})

	JustBeforeEach(func() {
		Expect(watcher.consumedRequests).Should(BeEmpty())

		go func() {
			defer GinkgoRecover()
			framework.ExpectNoError(grw.Watch(requests))
		}()
	})

	AfterEach(func() {
		close(requests)
		Expect(requests).Should(BeClosed())
	})

	Context("when sending new reconcile request", func() {
		It("watcher processes the request", func() {
			newRequest := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "ns", Name: "string"}}
			requests <- newRequest
			time.Sleep(waitPerRequest)
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(newRequest.String(), 1))
		})
	})

	Context("when sending an empty reconcile request", func() {
		It("watcher errors processing the request", func() {
			newRequest := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "", Name: ""}}
			requests <- newRequest
			time.Sleep(waitPerRequest)
			Expect(watcher.consumedRequests).ShouldNot(HaveKeyWithValue(newRequest.String(), 1))
		})
	})

	Context("when sending repeated requests", func() {
		It("watcher processes them one by one", func() {
			request := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "ns", Name: "first"}}
			By("sending first request")
			requests <- request
			// so that Watch goroutine can atleast execute
			time.Sleep(waitPerRequest / 4)
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request.String(), 1))
			By("sending second request")
			requests <- request
			By("check active channels")
			Expect(grw.activeChannels).Should(HaveLen(1))
			By("check request consumed")
			// first request hasn't finished processing
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request.String(), 1))
			time.Sleep(waitPerRequest)
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request.String(), 2))
		})
	})

	Context("when sending distinct requests", func() {
		It("watcher processes them in parallel", func() {
			request1 := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "ns", Name: "first"}}
			request2 := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "ns", Name: "second"}}
			By("sending first request")
			requests <- request1
			// so that Watch goroutine can atleast execute
			time.Sleep(waitPerRequest / 4)
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request1.String(), 1))
			By("sending second request")
			requests <- request2
			time.Sleep(waitPerRequest / 4)
			By("check active channels")
			Expect(grw.activeChannels).Should(HaveLen(2))
			By("check request consumed")
			// second request has already started processing
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request1.String(), 1))
			Expect(watcher.consumedRequests).Should(HaveKeyWithValue(request2.String(), 1))
		})
	})
})
