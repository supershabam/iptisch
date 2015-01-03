package iptisch

type Watcher interface {
	Close()
	Err() error
	Watch(keys []string) <-chan map[string][]string
}
