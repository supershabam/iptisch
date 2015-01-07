package iptisch

import (
	"io/ioutil"
	"os/exec"
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
		out, err := exec.Command("/sbin/iptables-save").Output()
		if err != nil {
			return err
		}
		iptisch := strings.Split(next, "\n")
		filtered := []string{}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		filterTable := false
		for scanner.Scan() {
			line := scanner.Text()
			// state of which table we're in. Need to know if we're in filter
			if strings.HasPrefix(line, "*") {
				if strings.HasPrefix(line, "*filter") {
					filterTable = true
				} else {
					filterTable = false
				}
			}
			// if we're modifying an IPTISCH chain (first 3 characters should be "-(A|I) "
			// we want to keep lines that route traffic into the IPTISCH chain e.g. -A INPUT -j IPTISCH_INPUT
			// TODO do this better
			if filterTable && strings.HasPrefix(line[3:], "IPTISCH") {
				continue
			}
			// inject iptisch content right before commmit
			if filterTable && strings.Contains(line, "COMMIT") {
				filtered = append(filtered, iptisch...)
			}
			filtered = append(filtered, line)
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		cmd := exec.Command("/sbin/iptables-restore")
		cmd.Stdin = strings.NewReader(strings.Join(filtered, "\n"))
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	if err := watcher.Err(); err != nil {
		return err
	}
	return nil
}
