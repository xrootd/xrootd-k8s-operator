package reconciler

import (
	"sync"

	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

const defaultChannelCapacity = 12

type WatchReconciler interface {
	RefreshWatch(request reconcile.Request) error
	GetWatchers(instance resource) watch.Watchers
	StartWatching() error
}

var _ inject.Stoppable = &WatchManager{}

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

func (wm *WatchManager) AddWatchers(watchers ...watch.Watcher) {
	wm.watchers = wm.watchers.AddWatchers(watchers...)
}

func (wm *WatchManager) InjectStopChannel(stop <-chan struct{}) error {
	wm.stop = stop
	return nil
}

func (wm *WatchManager) GetWatchers(instance resource) watch.Watchers {
	return wm.watchers
}

func (wm *WatchManager) RefreshWatch(request reconcile.Request) error {
	wm.destLock.Lock()
	wm.destLock.Unlock()

	for _, dst := range wm.dest {
		dst <- request
	}
	return nil
}

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
	dst := make(chan reconcile.Request, defaultChannelCapacity)
	go watcher.Watch(dst)
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

func NewWatchManager(stop <-chan struct{}) *WatchManager {
	initialSize := 0
	return &WatchManager{
		watchers: watch.Watchers(make([]watch.Watcher, initialSize)),
		stop:     stop,
		once:     sync.Once{},
		dest:     make([]chan reconcile.Request, initialSize),
		destLock: sync.Mutex{},
	}
}
