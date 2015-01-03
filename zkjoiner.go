package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type ZKJoiner struct {
	Conn        *zk.Conn
	Memberships string
	Root        string

	done chan struct{}
}

func (j *ZKJoiner) Close() {
	close(j.done)
}

func (j *ZKJoiner) Join() error {
	j.done = make(chan struct{})
	for _, membership := range strings.Split(j.Memberships, ",") {
		parts := strings.Split(membership, "+=")
		group := parts[0]
		ip := parts[1]
		znode := fmt.Sprintf("%s%s/%s-", j.Root, group, ip)
		_, err := j.Conn.Create(znode, []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	<-j.done
	return nil
}
