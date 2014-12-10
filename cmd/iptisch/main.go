package main

import (
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/supershabam/iptisch"
)

var (
	output  string        // file to write on update
	restore string        // path to iptables-restore command
	period  time.Duration // pause between polls
	dsn     string        // database
)

// failure: database not available
// - avoid trampling herd
// - last successful run should be loaded in iptables
//
// failure: generated firewall rules are invalid
// - log failure, do not save file
//

var (
	template = flag.String("template", "", "template to read")
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
		Keys:   []string{"key1", "key2"},
		Period: time.Second * 2,
	}
	for variables := range w.Watch() {
		log.Printf("%s\n\n", t.Execute(variables))
	}
	if w.Err != nil {
		log.Fatal(w.Err)
	}
}
