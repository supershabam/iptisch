package iptisch

import (
	"fmt"
	"math/rand"
	"time"
)

// A Watcher provides the variables your template needs. This one is bullshit
// because the implementation doesn't matter.
type Watcher struct {
	Err    error
	Keys   []string
	Period time.Duration
}

func (w *Watcher) Watch() <-chan Variables {
	out := make(chan Variables)
	go func() {
		defer close(out)
		for {
			out <- gen(w.Keys)
			time.Sleep(w.Period)
		}
	}()
	return out
}

func gen(keys []string) Variables {
	v := Variables{}
	for _, key := range keys {
		for _, n := range rand.Perm(rand.Intn(20)) {
			value := fmt.Sprintf("10.0.0.%d", n)
			if _, ok := v[key]; !ok {
				v[key] = []string{}
			}
			v[key] = append(v[key], value)
		}
	}
	return v
}
