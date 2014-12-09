package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/supershabam/iptisch"
)

var (
	servers  = flag.String("servers", "", "zookeeper servers to connect to")
	template = flag.String("template", "", "path to template file")
	znode    = flag.String("znode", "", "zookeeper znode to watch for data")
)

func main() {
	flag.Parse()
	templateData, err := ioutil.ReadFile(*template)
	if err != nil {
		log.Fatal(err)
	}
	t := iptisch.Template{
		Text: string(templateData),
	}
	w := iptisch.Watcher{
		Servers: strings.Split(*servers, ","),
		ZNode:   *znode,
	}
	for variables := range w.Watch() {
		log.Printf("%s\n\n", t.Execute(variables))
	}
	if w.Err != nil {
		log.Fatal(w.Err)
	}
}
