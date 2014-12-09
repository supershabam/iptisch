package iptisch

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type Values interface {
	Err() error
	Watch() <-chan map[string][]byte
}

type ZKValues struct {
	base map[string][]byte
	conn *zk.Conn
	err  error
	done chan struct{}
	keys []string
}

func NewZKValues(servers, keys []string) (*ZKValues, error) {
	zkValues := &ZKValues{
		base: map[string][]byte{},
		done: make(chan struct{}),
		keys: keys,
	}

	conn, _, err := zk.Connect(servers, time.Minute)
	if err != nil {
		return nil, err
	}
	zkValues.conn = conn

	for _, key := range keys {
		value, _, err := conn.Get(key)
		if err != nil {
			return nil, err
		}
		zkValues.base[key] = value
	}

	return zkValues, nil
}

func (ZKValues v) Err() error {
	return v.err
}

func (ZKValues *v) Watch() <-chan map[string][]byte {
	out := make(chan map[string][]byte)
	go func() {
		defer close(out)
		// send pre-loaded values and then updates as they come
		out <- v.values()
	}()
	return out
}

func (ZKValues *v) watchKey(key string) error {
	for {
		value, _, event, err := v.conn.GetW(key)
		if err != nil {
			return err
		}
		select {
		case <-v.done:
			return nil
		case <-event:
			// loop
		}
	}
}

func (ZKValues v) values() map[string][]byte {
	v := map[string][]byte{}
	for key, value := range v.base {
		v[key] = value
	}
	return v
}
