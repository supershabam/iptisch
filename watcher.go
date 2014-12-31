package iptisch

import (
	"fmt"
	"math/rand"
	"time"
)

type Group struct {
	Name string
	IPs  []string
}

type GroupWatcher interface {
	Close()
	Err() error
	Watch() <-chan Group
}

type Groups map[string]Group

type GroupsWatcher interface {
	Close()
	Err() error
	Watch() <-chan Groups
}
