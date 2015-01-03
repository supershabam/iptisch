package main

import (
	"flag"
	"log"
	"sync"

	"github.com/supershabam/iptisch"
)

var (
	command     = flag.String("command", "", "command to execute each time template is compiled")
	memberships = flag.String("memberships", "", "comma separated list of group+=ip")
	servers     = flag.String("servers", "", "comma separated list of zookeeper addresses")
	template    = flag.String("template", "", "template to execute")
)

func main() {
	wg := sync.WaitGroup{}
	done := make(chan struct{})
	flag.Parse()
	factory := iptisch.Factory{
		Servers: *servers,
	}
	if len(*memberships) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			joiner, err := factory.Joiner(*memberships)
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				<-done
				joiner.Close()
			}()
			if err := joiner.Join(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	if len(*template) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			watcher, err := factory.Watcher()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				<-done
				watcher.Close()
			}()
			if err := iptisch.Run(watcher, *template, *command); err != nil {
				log.Fatal(err)
			}
		}()
	}
	wg.Wait()
}
