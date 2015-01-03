package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
)

// A ZKChildWatcher creates a channel of the current state of the children nodes
// in zookeeper that is updated as the children change.
type ZKChildWatcher struct {
	Conn  *zk.Conn
	Group string
	Root  string

	done chan struct{}
	err  error
}

// Close will close a currently executing watch. It is an error to call this function without
// having first called Watch.
func (cw *ZKChildWatcher) Close() {
	close(cw.done)
}

// Err is checked after the watch range completes to see if there was an error.
func (cw ZKChildWatcher) Err() error {
	return cw.err
}

// Watch returns a channel of child keys as they are updates in zookeeper, the first value
// on the channel is the current state of the children in zookeeper, and subsequent values
// are when the children are updated. If the group doesn't exist, then an error is set and
// the channel is closed.
// It is an error to call this function more than once.
func (cw *ZKChildWatcher) Watch() <-chan []string {
	cw.done = make(chan struct{})
	out := make(chan []string)
	go func() {
		defer close(out)
		zpath := fmt.Sprintf("%s%s", cw.Root, cw.Group)
		for {
			children, _, eventCh, err := cw.Conn.ChildrenW(zpath)
			if err != nil {
				cw.err = err
				return
			}
			out <- children
			select {
			case <-cw.done:
				return
			case event := <-eventCh:
				if event.Type == zk.EventNodeChildrenChanged {
					continue
				}
				cw.err = fmt.Errorf("unhandled zk event type: %s", event)
				return
			}
		}
	}()
	return out
}
