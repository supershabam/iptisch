package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/supershabam/iptisch"
)

var (
	template  = flag.String("template", "", "template file")
	variables = flag.String("variables", "", "variables json file")
)

func main() {
	flag.Parse()

	templateBytes, err := ioutil.ReadFile(*template)
	if err != nil {
		panic(fmt.Errorf("could not read template file: %s", *template))
	}

	variablesBytes, err := ioutil.ReadFile(*variables)
	if err != nil {
		panic(fmt.Errorf("could not read variables file: %s", *variables))
	}

	t := iptisch.Template{
		Text: string(templateBytes),
	}
	v := iptisch.Variables{}
	err = json.Unmarshal(variablesBytes, &v)
	if err != nil {
		panic(err)
	}

	t.Execute(os.Stdout, v)
}
