package watch

import (
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type watchFunc func(reconcile.Request) error

// GroupedRequestWatcher groups reconcile.Request based on NamespacedName and sends distinct ones
// to their goroutines.
type GroupedRequestWatcher struct {
	activeChannels map[string]chan<- reconcile.Request
	subWatcher     Watcher
}

var _ Watcher = &GroupedRequestWatcher{}

// Watch implements watch.Watcher
func (grw *GroupedRequestWatcher) Watch(requests <-chan reconcile.Request) error {
	logger := log.WithName("GroupedRequestWatcher.Watch")
	for request := range requests {
		logger.Info("Refreshing watch...", "request", request)
		key := request.String()
		channel := grw.getDistinctRequestChannel(key)
		if len(channel) < cap(channel) {
			channel <- request
		}
	}
	grw.doCleanup()
	return nil
}

func (grw *GroupedRequestWatcher) getDistinctRequestChannel(key string) chan<- reconcile.Request {
	channel, ok := grw.activeChannels[key]
	if !ok {
		channel = grw.startNewRequestChannel(key)
	}
	return channel
}

func (grw *GroupedRequestWatcher) startNewRequestChannel(key string) chan<- reconcile.Request {
	channel := make(chan reconcile.Request, 1)
	go grw.subWatcher.Watch(channel)

	grw.activeChannels[key] = channel

	return channel
}

func (grw *GroupedRequestWatcher) doCleanup() {
	for _, channel := range grw.activeChannels {
		close(channel)
	}
}

func NewGroupedRequestWatcher(subWatcher Watcher) Watcher {
	return &GroupedRequestWatcher{
		subWatcher:     subWatcher,
		activeChannels: make(map[string]chan<- reconcile.Request),
	}
}
