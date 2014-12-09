package iptisch

import "strings"

type Template struct {
	Text string
}

func (t Template) Execute(v map[string][]string) string {
	out := []string{}
	for _, line := range strings.Split(t.Text, "\n") {
		for _, expansion := range expand(line, v) {
			out = append(out, expansion)
		}
	}
	return strings.Join(out, "\n")
}

// simple template where a pipe becomes a special character and can't be escaped
// -A MYSQL -s |staging -j ACCEPT
// -A MYSQL -s |blacklist -j DROP
func expand(line string, v map[string][]string) (expansions []string) {
	// base case
	if !strings.Contains(line, "|") {
		expansions = append(expansions, line)
		return
	}
	index := strings.Index(line, "|")
	key := strings.SplitN(line[index:], " ", 2)[0]
	name := key[1:] // cut off the |
	values := v[name]
	for _, value := range values {
		replaced := strings.Replace(line, key, value, 1)
		// since multiple variables can exist in a line
		expansions = append(expansions, expand(replaced, v)...)
	}
	return
}
