mailfull-go
===========

A management tool for virtual domain email for Postfix and Dovecot written in Go.

[![GoDoc](https://godoc.org/github.com/directorz/mailfull-go?status.svg)](https://godoc.org/github.com/directorz/mailfull-go)

Features
--------

- You can use both virtual users and system users.
- Mailfull does not involve in delivery processes of the MTA, is only to generate configuration databases.
- You do not need to restart Postfix/Dovecot to apply configuration databases.
- The received email can be passed to the programs.

Installation
------------

### go get

Installed in `$GOPATH/bin`

```
$ go get github.com/directorz/mailfull-go/cmd/mailfull
```

Quick Start
-----------

Create a new user for Mailfull.

```
# useradd -r -s /bin/bash mailfull
# su - mailfull
```

Initialize a directory as a Mailfull repository.

```
$ mkdir /path/to/repo && cd /path/to/repo
$ mailfull init
```

Generate configurations for Postfix and Dovecot. (Edit as needed.)

```
$ mailfull genconfig postfix > /etc/postfix/main.cf
$ mailfull genconfig dovecot > /etc/dovecot/dovecot.conf
```

Start Postfix and Dovecot.

```
# systemctl start postfix.service
# systemctl start dovecot.service
```

Add a new domain and user.

```
# cd /path/to/repo

# mailfull domainadd example.com
# mailfull useradd hoge@example.com
# mailfull userpasswd hoge@example.com
```

Enjoy!

More info
---------

See [documentation](doc/README.md)
