iptisch
=======

[EXPERIMENT] reactive iptable rule template

## goal

iptables rules that have "security groups" where ip addresses can dynamically change membership with a group

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

### data model

Groups - named collections of ip addresses
IP Address - a single ipv4 ip address

A group is specified as a normal znode created under root. If a client requests a group that does not exit, that group znode is created.

If a client tries to register itself as a member of a group that does not exists, that group znode is created.

A ephemeral znode is created for an ip address under the group znode. The key of the znode is a generated prefix followed by the ip address. e.g. `aT23f-1.2.3.4`

Since a client can restart (or crash) it's ephemeral ip address znode under the group can still be present in the group when it tries to create a node with its ip address in the group. This is why we prefix the key.

We make the IP Address key contain the ip address because getting the keys of the children of the group is a single operation. Getting the children, and then reading the values of all the children would require N operations. So, the actual value of these IP znodes is "".

#### example tree

```txt
.
├── blacklist
│   ├── aTfEG-1.2.3.4
│   ├── s7GbA-1.2.3.4
│   ├── 8GvJm-4.3.2.1
├── frontend
│   ├── aTfEG-3.3.3.3
│   ├── s7GbA-4.4.4.4
│   ├── 8GvJm-5.5.5.5
```
