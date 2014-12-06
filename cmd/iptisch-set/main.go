package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	ZKAddr = "104.131.40.109:2181"
)

func main() {
	conn, _, err := zk.Connect([]string{ZKAddr}, time.Minute)
	if err != nil {
		panic(err)
	}
	exists, _, err := conn.Exists("/wut")
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err = conn.Create("/wut", []byte("nil"), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			panic(err)
		}
	}
	b, s, err := conn.Get("/wut")
	if err != nil {
		panic(err)
	}
	fmt.Printf("pulled: %s@%d\n", b, s.Version)
	_, err = conn.Set("/wut", []byte("next"), s.Version)
	if err != nil {
		panic(err)
	}
}
