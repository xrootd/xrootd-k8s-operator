package reconciler

import (
	"sync"

	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

const defaultChannelCapacity = 12

var log = logf.Log.WithName("reconciler")

// WatchReconciler provides control to watch lifecycle and allows refreshing
// watches on reconciliation request.
type WatchReconciler interface {
	RefreshWatch(request reconcile.Request) error
	GetWatchers(instance resource) watch.Watchers
	StartWatching() error
}

var _ inject.Stoppable = &WatchManager{}

// WatchManager manages a set of watchers, each of them running in parallel goroutines.
// It also supports graceful shutdown of watchers by closing all the destination
// channels whenever the injected stop channel is triggered.
type WatchManager struct {
	// once ensures the event distribution goroutine will be performed only once
	once sync.Once
	// stop is to end ongoing goroutine, and close the channels
	stop     <-chan struct{}
	watchers watch.Watchers
	// dest is the destination channels of the added watchers
	dest []chan reconcile.Request
	// destLock is to ensure the destination channels are safely added/removed
	destLock sync.Mutex
}

// AddWatchers adds a requested set of watchers to be managed
func (wm *WatchManager) AddWatchers(watchers ...watch.Watcher) {
	wm.watchers = wm.watchers.AddWatchers(watchers...)
}

// InjectStopChannel sets the stop channel, which when triggered would
// stop ongoing goroutines and close all channels
func (wm *WatchManager) InjectStopChannel(stop <-chan struct{}) error {
	wm.stop = stop
	return nil
}

// GetWatchers returns the set of managed watchers
func (wm *WatchManager) GetWatchers(instance resource) watch.Watchers {
	return wm.watchers
}

// RefreshWatch refreshes the managed-watchers with the
// given reconciliation request.
func (wm *WatchManager) RefreshWatch(request reconcile.Request) error {
	wm.destLock.Lock()
	defer wm.destLock.Unlock()

	for _, dst := range wm.dest {
		dst <- request
	}
	return nil
}

// StartWatching starts up the required goroutines of watchers, which
// either terminates due to error in Watch() or when all channels are closed.
func (wm *WatchManager) StartWatching() error {
	wm.once.Do(func() {
		go wm.syncLoop()
	})

	wm.destLock.Lock()
	defer wm.destLock.Unlock()

	for _, watcher := range wm.watchers {
		wm.doWatch(watcher)
	}

	return nil
}

func (wm *WatchManager) doWatch(watcher watch.Watcher) {
	logger := log.WithName("watcher").WithName("doWatch")
	dst := make(chan reconcile.Request, defaultChannelCapacity)
	go func() {
		if err := watcher.Watch(dst); err != nil {
			logger.Error(err, "Watcher errored!")
		}
		logger.Info("Watcher finished watching...")
	}()
	wm.dest = append(wm.dest, dst)
}

func (wm *WatchManager) doStop() {
	wm.destLock.Lock()
	defer wm.destLock.Unlock()

	for _, dst := range wm.dest {
		close(dst)
	}
}

func (wm *WatchManager) syncLoop() {
	for {
		<-wm.stop
		// Close destination channels
		wm.doStop()
		return
	}
}

// NewWatchManager creates a new WatchManager with empty watchers
// and given stop channel
func NewWatchManager(stop <-chan struct{}) *WatchManager {
	initialSize := 0
	return &WatchManager{
		watchers: watch.Watchers(make([]watch.Watcher, initialSize)),
		stop:     stop,
		dest:     make([]chan reconcile.Request, initialSize),
	}
}
