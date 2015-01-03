package iptisch

import (
	"sync"

	"github.com/samuel/go-zookeeper/zk"
)

type ZKWatcher struct {
	Conn *zk.Conn
	Root string

	err  error
	done chan struct{}
}

func (w *ZKWatcher) Close() {
	close(w.done)
}

// Err returns nil if there was no error after Watch closes. Otherwise, it returns the
// error which caused Watch to close.
func (w ZKWatcher) Err() error {
	return w.err
}

// Watch returns a map of group to its children key values. It will only send a map once
// a value is retreived for each group (no partial map). Then, it sends a full map immediately
// after any single group has its values updated.
// If any error occurs along the way, the channel is closed and Err() will return the error.
// If any zookeeper node doesn't exist, then it is an error
func (w *ZKWatcher) Watch(keys []string) <-chan map[string][]string {
	watchers := []*ZKChildWatcher{}
	for _, key := range keys {
		watchers = append(watchers, &ZKChildWatcher{
			Conn:  w.Conn,
			Group: key,
			Root:  w.Root,
		})
	}
	w.done = make(chan struct{})
	out := make(chan map[string][]string)
	wg := sync.WaitGroup{}
	l := sync.Mutex{}
	m := map[string][]string{}

	wg.Add(len(watchers))
	for _, watcher := range watchers {
		go func(watcher *ZKChildWatcher) {
			defer wg.Done()
			cch := watcher.Watch()
		Drain:
			for {
				select {
				case <-w.done:
					w.done = nil    // don't re-enter this case
					watcher.Close() // will cause cch to close when drained
				case children, next := <-cch:
					if !next {
						break Drain
					}
					l.Lock()
					m[watcher.Group] = children
					count := len(m)
					l.Unlock()
					if count == len(watchers) {
						out <- m
					}
				}
			}
			l.Lock()
			if err := w.Err(); err != nil && w.err == nil {
				w.err = err
				watcher.Close()
			}
			l.Unlock()
		}(watcher)
	}

	go func() {
		wg.Wait()
		defer close(out)
	}()
	return out
}
