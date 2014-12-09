package iptisch

import (
	"encoding/json"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Watcher struct {
	Err     error
	Servers []string
	ZNode   string
}

func (w *Watcher) Watch() <-chan Variables {
	out := make(chan Variables)
	go func() {
		defer close(out)
		conn, _, err := zk.Connect(w.Servers, time.Minute)
		if err != nil {
			w.Err = err
			return
		}
		for {
			data, _, event, err := conn.GetW(w.ZNode)
			if err != nil {
				w.Err = err
				return
			}
			var variables Variables
			err = json.Unmarshal(data, &variables)
			if err != nil {
				w.Err = err
				return
			}
			out <- variables
			<-event // reloop when data changes
		}
	}()
	return out
}
