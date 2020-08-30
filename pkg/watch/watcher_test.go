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
	"fmt"
	"time"

	"github.com/onsi/ginkgo"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type TestWatcher struct {
	consumedRequests map[string]int
	DelayPerRequest  time.Duration
}

var _ Watcher = TestWatcher{}

func (tw TestWatcher) Watch(requests <-chan reconcile.Request) error {
	defer ginkgo.GinkgoRecover()
	for request := range requests {
		key := request.String()
		if len(request.Name) == 0 {
			return fmt.Errorf("empty name in request not allowed")
		}
		tw.consumedRequests[key]++
		time.Sleep(tw.DelayPerRequest)
	}
	return nil
}

func newTestWatcher(delayPerRequest time.Duration) TestWatcher {
	return TestWatcher{
		DelayPerRequest:  delayPerRequest,
		consumedRequests: make(map[string]int),
	}
}
