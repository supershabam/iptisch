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
*filter
-A INPUT -s |frontend -j ACCEPT
-A OUTPUT -d |frontend -j ACCEPT
```

Then, on your front-end servers run the service `iptisch -membership="frontend+=$ip_address"`

Now your font-end servers will be allowed to talk to your backend server.

