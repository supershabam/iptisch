package main

import (
	"flag"
	"log"

	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/supershabam/iptisch"
	"strings"
	"time"
)

var (
	memberships = flag.String("memberships", "", "comma separated list of group+=ip")
	servers     = flag.String("servers", "", "comma separated list of zookeeper addresses")
	root        = flag.String("root", "/", "zookeeper root (for namespacing)")
)

func main() {
	flag.Parse()

	conn, _, err := zk.Connect(strings.Split(*servers, ","), time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	gw := iptisch.GroupsWatcher{
		Watchers: []*iptisch.ChildWatcher{
			&iptisch.ChildWatcher{
				Conn:  conn,
				Group: "iptisch",
				Root:  *root,
			},
			&iptisch.ChildWatcher{
				Conn:  conn,
				Group: "test",
				Root:  *root,
			},
		},
	}
	err = iptisch.WriteMemberships(conn, *root, *memberships)
	if err != nil {
		log.Fatal(err)
	}

	for group := range gw.Watch() {
		fmt.Printf("%+v\n", group)
	}

	if gw.Err() != nil {
		log.Fatal(gw.Err())
	}
}
