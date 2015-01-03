package iptisch

type Joiner interface {
	Close()
	Join() error
}
