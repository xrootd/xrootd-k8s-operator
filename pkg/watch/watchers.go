package watch

// Watchers is a utility struct to deal with multiple watchers
type Watchers []Watcher

// ToSlice converts Watchers to primitive array - []Watcher
func (ws Watchers) ToSlice() []Watcher {
	return []Watcher(ws)
}

// AddWatchers adds passed watchers
func (ws *Watchers) AddWatchers(watcher ...Watcher) Watchers {
	return append(ws.ToSlice(), watcher...)
}
