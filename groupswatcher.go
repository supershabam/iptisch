package iptisch

import (
	"sync"
)

// A GroupsWatcher observes multiple zookeeper nodes and their children.
type GroupsWatcher struct {
	Watchers []*ChildWatcher

	err  error
	done chan struct{}
}

// Close will cause watch to close. It is an error to call this function without first
// calling Watch
func (gw *GroupsWatcher) Close() {
	close(gw.done)
}

// Err returns nil if there was no error after Watch closes. Otherwise, it returns the
// error which caused Watch to close.
func (gw GroupsWatcher) Err() error {
	return gw.err
}

// Watch returns a map of group to its children key values. It will only send a map once
// a value is retreived for each group (no partial map). Then, it sends a full map immediately
// after any single group has its values updated.
// If any error occurs along the way, the channel is closed and Err() will return the error.
// If any zookeeper node doesn't exist, then it is an error
func (gw *GroupsWatcher) Watch() <-chan map[string][]string {
	gw.done = make(chan struct{})
	out := make(chan map[string][]string)
	wg := sync.WaitGroup{}
	l := sync.Mutex{}
	m := map[string][]string{}

	wg.Add(len(gw.Watchers))
	for _, w := range gw.Watchers {
		go func(w *ChildWatcher) {
			defer wg.Done()
			cch := w.Watch()
		Drain:
			for {
				select {
				case <-gw.done:
					gw.done = nil // don't re-enter this case
					w.Close()     // will cause cch to close when drained
				case children, next := <-cch:
					if !next {
						break Drain
					}
					l.Lock()
					m[w.Group] = children
					count := len(m)
					l.Unlock()
					if count == len(gw.Watchers) {
						out <- m
					}
				}
			}
			l.Lock()
			if err := w.Err(); err != nil && gw.err == nil {
				gw.err = err
				gw.Close()
			}
			l.Unlock()
		}(w)
	}

	go func() {
		wg.Wait()
		defer close(out)
	}()
	return out
}
