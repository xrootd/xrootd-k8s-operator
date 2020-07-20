package watch

import "sigs.k8s.io/controller-runtime/pkg/reconcile"

type channelRequest chan reconcile.Request

type watchFunc func(reconcile.Request)

// GroupedRequestWatcher groups reconcile.Request based on NamespacedName and sends distinct ones
// to their goroutines.
type GroupedRequestWatcher struct {
	activeChannels map[string]channelRequest
	Func           watchFunc
}

var _ Watcher = GroupedRequestWatcher{}

// Watch implements watch.Watcher
func (grw GroupedRequestWatcher) Watch(request reconcile.Request) error {
	key := request.String()
	channel, ok := grw.activeChannels[key]
	if !ok {
		channel = make(channelRequest, 1)
		grw.activeChannels[key] = channel
		go func() {
			for req := range channel {
				grw.Func(req)
			}
		}()
	}
	if len(channel) < cap(channel) {
		channel <- request
	}
	return nil
}

func NewGroupedRequestWatcher(function watchFunc, stop <-chan struct{}) GroupedRequestWatcher {
	return GroupedRequestWatcher{
		Func:           function,
		activeChannels: make(map[string]channelRequest),
	}
}
