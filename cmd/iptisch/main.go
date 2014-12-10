package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/supershabam/iptisch"
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
	out      = flag.String("out", "", "output file to write")
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
		Keys:   t.Keys(),
		Period: time.Second * 2,
	}
	for variables := range w.Watch() {
		ruleset := bytes.NewBufferString(t.Execute(variables))
		err := ioutil.WriteFile(*out, ruleset.Bytes(), 0664)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote ruleset")
	}
	if w.Err != nil {
		log.Fatal(w.Err)
	}
}
