package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
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

func Run(conn *zk.Conn, root, path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	t := Template{string(raw)}
	watchers := []*ChildWatcher{}
	for _, group := range t.Keys() {
		watchers = append(watchers, &ChildWatcher{
			Conn:  conn,
			Group: group,
			Root:  root,
		})
	}
	gw := GroupsWatcher{
		Watchers: watchers,
	}
	last := ""
	for m := range gw.Watch() {
		next := t.Execute(ipMap(m))
		// dedupe if we generate the same result
		if next == last {
			continue
		}
		last = next
		fmt.Println(next)
	}
	if err := gw.Err(); err != nil {
		return err
	}
	return nil
}
