package iptisch

import (
	"fmt"
	"math/rand"
	"time"
)

type IP string

type ZKGroupWatcher struct {
	Conn  *zk.Conn
	Group string

	done chan struct{}
	err  error
}

func childrenToIPs(children []string) []IP {
	ips := make([]IP)
	for _, child := range children {
		// TODO remove prefix
		// TODO dedupe
		ips = append(ips, IP(child))
	}
	return ips
}

func (gw *ZKGroupWatcher) Close() {
	close(gw.done)
}

func (gw *ZKGroupWatcher) Watch() <-chan Group {
	gw.done = make(chan struct{})
	out := make(chan []IP)
	go func() {
		defer close(out)
		for {
			children, _, eventCh, err := gw.Conn.ChildrenW(gw.Group)
			if err != nil {
				gw.Err = err
				return
			}
			out <- childrenToIPs(children)
			select {
			case <-done:
				return
			case event := <-eventCh:
				if event != zk.EventNodeChildrenChanged {
					gw.Err = fmt.Errorf("unhandled node event on group: %s = %s", gw.Group, event)
					return
				}
			}
		}
	}()
	return out
}
