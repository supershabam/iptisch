package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
)

func Run(conn *zk.Conn, root, path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	t := Template{string(raw)}
	watchers := []*ChildWatcher{}
	for _, group := range t.Keys() {
		watchers = append(watchers, &ChildWatcher{
			Conn:  conn,
			Group: group,
			Root:  root,
		})
	}
	gw := GroupsWatcher{
		Watchers: watchers,
	}
	for m := range gw.Watch() {
		fmt.Println(t.Execute(m))
	}
	if err := gw.Err(); err != nil {
		return err
	}
	return nil
}
