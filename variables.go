package iptisch

import "fmt"

type Variables struct {
}

func (Variables v) Get(key string) (values []string) {
	for i := 0; i < len(key); i++ {
		values = append(values, fmt.Sprintf("%s%d", key, i))
	}
	return
}
