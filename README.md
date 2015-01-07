iptisch
=======

[ALPHA] reactive iptables

[![Build Status](https://travis-ci.org/supershabam/iptisch.svg?branch=master)](https://travis-ci.org/supershabam/iptisch)

*your filewall should change as your infrastructure changes*

example: securing a database server
-----------------------------------

You have 1 database server and multiple front-end servers that come and go as your traffic changes.

From your database server run the service `iptisch -template="/etc/iptisch/rules"`

Where `/etc/iptisch/rules` contains something like this

```
```

## using the binary

### specifying group membership

You can run the long-running iptisch binary to specify that your machine is part of the `whitelist` and `anothergroup` groups.

`iptisch -servers="$zookeeper_servers" -membership="whitelist+=$ip_address,anothergroup+=$ip_address"`

### being notified on group change

This long-running iptisch binary will execute command each time the template is executed with membership values from zookeeper.

`iptisch -servers="$zookeeper_servers" -template="/path/to/template.eit" -command="/custom/dosomething.sh"`

### combining!

You can specify membership and execute a template in one iptisch binary

`iptisch -servers="$zookeeper_servers" -template="/path/to/template.eit" -command="/custom/dosomething.sh" -membership="whitelist+=$ip_address,anothergroup+=$ip_address"`

## templating

The format of an iptables rule file is known. It would be cumbersome to write iteration logic into the template to list all ips of a group. So, let's try and invent our own simple templating language.

### expanding iptables `.eit`

The expanding iptables template is a very simple templating language: `|` is a special character which causes the current line to explode with members of the group specified after the `|`.

#### blacklist example

```eit
# blacklist rule (text without pipes are preserved)
-A INPUT -s |blacklist -j DROP
# of course, you can also write non-expanding rules for non-changing rules
-A INPUT -s 10.0.0.0/8 -j ACCEPT
```

output given `blacklist=['1.2.3.4/32', '4.3.2.1/16']`

```txt
# blacklist rule (text without pipes are preserved)
-A INPUT -s 1.2.3.4 -j DROP
-A INPUT -s 4.3.2.1 -j DROP
# of course, you can also write non-expanding rules for non-changing rules
-A INPUT -s 10.0.0.0/8 -j ACCEPT
```

output given `blacklist=[]`

```txt
# blacklist rule (text without pipes are preserved)
# of course, you can also write non-expanding rules for non-changing rules
-A INPUT -s 10.0.0.0/8 -j ACCEPT
```

## using zookeeper as a backend

If a computer wants its IP address to be part of a group, it needs to connect to zookeeper and put it's IP address into the group. If that server goes away, that membership should also go away. Zookeeper's ephemeral znodes provides this functionality without needing to write custom heartbeating logic.

A iptisch client needs to be notified when the membership of a group changes and given the new membership. Zookeeper provides this with watches.

### data structure

ZooKeeper is a tree structure. Groups are the first nodes, and ip addresses are the children of groups. ZooKeeper has a built-in sequential node to prevent collisions (two places where 1.2.3.4 is defined to be part of the blacklist group).

```txt
.
├── blacklist
│   ├── 1.2.3.4-00000001
│   ├── 1.2.3.4-00000002
│   ├── 4.3.2.1-00000001
├── frontend
│   ├── 3.3.3.3-00000001
│   ├── 4.4.4.4-00000001
│   ├── 5.5.5.5-00000001
```

