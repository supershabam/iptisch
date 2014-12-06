package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/supershabam/iptisch"
)

var (
	template = flag.String("template", "", "template file")
)

const (
	ZKAddr = "104.131.40.109:2181"
)

func main() {
	flag.Parse()

	templateBytes, err := ioutil.ReadFile(*template)
	if err != nil {
		panic(fmt.Errorf("could not read template file: %s", *template))
	}

	t := iptisch.Template{
		Text: string(templateBytes),
	}

	v := iptisch.Variables{
		Keys: []string{
			"/test",
		},
		Servers: []string{ZKAddr},
	}
	for m := range v.Watch() {
		t.Execute(os.Stdout, m)
	}
	if v.Err != nil {
		panic(v.Err)
	}
	fmt.Printf("done\n")
}
