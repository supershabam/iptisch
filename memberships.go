package iptisch

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

func WriteMemberships(conn *zk.Conn, root, memberships string) error {
	for _, membership := range strings.Split(memberships, ",") {
		parts := strings.Split(membership, "+=")
		group := parts[0]
		ip := parts[1]
		znode := fmt.Sprintf("%s%s/%s-", root, group, ip)
		_, err := conn.Create(znode, []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	return nil
}
