package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
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
		cmd := exec.Command("iptables-restore")
		cmd.Stdin = strings.NewReader(t.Execute(variables))
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("%s\n", out)
			log.Fatal(err)
		}
		log.Printf("wrote rules to iptables")
	}
	if w.Err != nil {
		log.Fatal(w.Err)
	}
}
