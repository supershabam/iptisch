package main

import (
	"fmt"

	"github.com/supershabam/iptisch"
)

const (
	ZKAddr = "104.131.40.109:2181"
)

func main() {
	v := iptisch.Variables{
		Keys: []string{
			"/test",
			"/wut",
		},
		Servers: []string{ZKAddr},
	}
	for m := range v.Watch() {
		fmt.Printf("%+v\n", m)
	}
	if v.Err != nil {
		panic(v.Err)
	}
	fmt.Printf("done\n")
}
