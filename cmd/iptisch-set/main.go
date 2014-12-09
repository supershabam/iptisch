package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	Servers = flag.String("servers", "", "zookeeper servers (comma separated) to connect to")
	ZNode   = flag.String("znode", "", "zookeeper znode path")
)

func main() {
	flag.Parse()
	conn, _, err := zk.Connect(strings.Split(*Servers, ","), time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	// try to create and set if not currently existing
	exists, _, err := conn.Exists(*ZNode)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		_, err = conn.Create(*ZNode, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// get and set and fail in race condition with other sets
	_, s, err := conn.Get(*ZNode)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Set(*ZNode, data, s.Version)
	if err != nil {
		log.Fatal(err)
	}
}
