package main

import (
	"flag"
	"log"

	"github.com/supershabam/iptisch"
)

var (
	servers = flag.String("servers", "", "comma separated list of zookeeper addresses")
)

func main() {
	flag.Parse()

	gw := iptisch.ZKGroupWatcher{
		Conn:  conn,
		Group: "test",
	}

	for group := range gw.Watch() {
		fmt.Printf("%+v\n", group)
	}

	if gw.Err() != nil {
		log.Fatal(err)
	}
}
