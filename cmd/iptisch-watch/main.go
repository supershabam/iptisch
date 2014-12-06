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
	for {
		b, s, c, err := conn.GetW("/test")
		if err != nil {
			panic(err)
		}
		fmt.Printf("has data: %s@%d\n", b, s.Version)
		<-c
	}
}
