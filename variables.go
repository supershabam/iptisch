package iptisch

import (
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Variables struct {
	Err     error
	Keys    []string
	Servers []string
}

func (v *Variables) Watch() <-chan map[string][]byte {
	out := make(chan map[string][]byte)
	values := map[string][]byte{}

	go func() {
		defer close(out)

		conn, _, err := zk.Connect(v.Servers, time.Minute)
		if err != nil {
			v.Err = err
			return
		}

		done := make(chan struct{})
		l := &sync.Mutex{}
		for _, key := range v.Keys {
			go func(key string) {
				var drainErr error = nil
				for value := range drain(done, conn, &drainErr, key) {
					l.Lock()
					values[key] = value
					// copy bytes so receiver has immutable copy of current state
					result := map[string][]byte{}
					for key, value := range values {
						result[key] = value
					}
					l.Unlock() // unlock before potentially blocking on send
					out <- result
				}
				if drainErr != nil {
					l.Lock()
					defer l.Unlock()
					if v.Err == nil {
						v.Err = drainErr
						close(done)
					}
				}
				return
			}(key)
		}
		<-done
	}()
	return out
}

func drain(done <-chan struct{}, conn *zk.Conn, err *error, key string) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		for {
			value, _, event, err := conn.GetW(key)
			if err != nil {
				return
			}
			out <- value
			select {
			case <-done:
				return
			case <-event:
			}
		}
	}()
	return out
}
