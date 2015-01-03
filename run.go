package iptisch

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func deSuffix(input []string) (output []string) {
	for _, i := range input {
		output = append(output, strings.Split(i, "-")[0])
	}
	return
}

func dedupe(input []string) (output []string) {
	seen := map[string]bool{}
	for _, i := range input {
		if seen[i] {
			continue
		}
		seen[i] = true
		output = append(output, i)
	}
	return
}

func ipMap(input map[string][]string) (output map[string][]string) {
	output = map[string][]string{}
	for group, values := range input {
		values = dedupe(deSuffix(values))
		sort.Strings(values)
		output[group] = values
	}
	return
}

func Run(watcher Watcher, template, command string) error {
	raw, err := ioutil.ReadFile(template)
	if err != nil {
		return err
	}
	t := Template{string(raw)}
	last := ""
	for m := range watcher.Watch(t.Keys()) {
		next := t.Execute(ipMap(m))
		// dedupe if we generate the same result
		if next == last {
			continue
		}
		last = next
		fmt.Println(next)
	}
	if err := watcher.Err(); err != nil {
		return err
	}
	return nil
}
