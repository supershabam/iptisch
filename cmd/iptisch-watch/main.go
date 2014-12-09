package main

import (
	"flag"
	"log"
	"strings"

	"github.com/supershabam/iptisch"
)

var (
	servers = flag.String("servers", "", "zookeeper servers to connect to")
	znode   = flag.String("znode", "", "zookeeper znode to watch for data")
)

func main() {
	flag.Parse()
	w := iptisch.Watcher{
		Servers: strings.Split(*servers, ","),
		ZNode:   *znode,
	}
	for variables := range w.Watch() {
		log.Printf("%+v", variables)
	}
	if w.Err != nil {
		log.Fatal(w.Err)
	}
}
