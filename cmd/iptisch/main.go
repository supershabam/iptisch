package main

import (
	"flag"
	"log"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/supershabam/iptisch"
	"strings"
	"time"
)

var (
	memberships = flag.String("memberships", "", "comma separated list of group+=ip")
	root        = flag.String("root", "/", "zookeeper root (for namespacing)")
	servers     = flag.String("servers", "", "comma separated list of zookeeper addresses")
	template    = flag.String("template", "", "template to execute")
)

func main() {
	flag.Parse()

	conn, _, err := zk.Connect(strings.Split(*servers, ","), time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	if len(*memberships) > 0 {
		err = iptisch.WriteMemberships(conn, *root, *memberships)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*template) > 0 {
		err = iptisch.Run(conn, *root, *template)
		if err != nil {
			log.Fatal(err)
		}
	}
	done := make(chan struct{})
	<-done // doesn't get done for now
}
