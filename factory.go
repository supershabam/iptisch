package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"sync"
	"time"
)

type zkinfo struct {
	conn *zk.Conn
	root string
}

func parseZkinfo(servers string) (*zkinfo, error) {
	servers = servers[len("zk://"):]
	// TODO parse root and duration out of servers string
	conn, _, err := zk.Connect(strings.Split(servers, ","), time.Minute)
	if err != nil {
		return nil, err
	}
	return &zkinfo{
		conn: conn,
		root: "/",
	}, nil
}

type Factory struct {
	Servers string

	m      sync.Mutex
	zkinfo *zkinfo
}

func (f *Factory) getzkinfo() (*zkinfo, error) {
	f.m.Lock()
	defer f.m.Unlock()
	if f.zkinfo == nil {
		zkinfo, err := parseZkinfo(f.Servers)
		if err != nil {
			return nil, err
		}
		f.zkinfo = zkinfo
	}
	return f.zkinfo, nil
}

func (f *Factory) Joiner(memberships string) (Joiner, error) {
	if strings.HasPrefix(f.Servers, "zk://") {
		zkinfo, err := f.getzkinfo()
		if err != nil {
			return nil, err
		}
		return &ZKJoiner{
			Conn:        zkinfo.conn,
			Memberships: memberships,
			Root:        zkinfo.root,
		}, nil
	}
	return nil, fmt.Errorf("nicht implementiert")
}

func (f *Factory) Watcher() (Watcher, error) {
	if strings.HasPrefix(f.Servers, "zk://") {
		zkinfo, err := f.getzkinfo()
		if err != nil {
			return nil, err
		}
		return &ZKWatcher{
			Conn: zkinfo.conn,
			Root: zkinfo.root,
		}, nil
	}
	return nil, fmt.Errorf("nicht gefunden")
}
