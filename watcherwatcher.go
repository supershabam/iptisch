package iptisch

import (
	"fmt"
	"math/rand"
	"time"
)

type WatcherWatcher struct {
	Conn     *zk.Conn
	Err      error
	Watchers []Watcher
}

func (ww *WatcherWatcher) Watch() <-chan map[string][]IP {
	out := make(chan map[string][]IP)
	state := map[string][]IP{}
	lock := &sync.Mutex{}
	go func() {
		defer close(out)
		for _, group := range w.Groups {
			go func(group string) {
			}(group)
		}
	}()
	return out
}
