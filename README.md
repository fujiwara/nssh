# nssh

nssh is a tool to execute command in parallel on multiple hosts and returns the aggregated result.

## Installation

### Homebrew

```console
$ brew install fujiwara/tap/nssh
```

### Binary package

[Releases](https://github.com/fujiwara/nssh/releases)

## Useage

```
$ nssh --help
Usage of nssh:
-i=false: handle stdin
-p=false: add hostname to line prefix
-t=[]: target hostname
-f=[]: target hostname from files
-v=false: show version
```

## Sample

```
$ nssh -t host-web-001 -t host-web-002 cat /var/log/nginx/access.log
12.34.567.890 - - [06/Oct/2015:15:52:42 +0900] "-" 400 0 "-" "-"
12.34.567.890 - - [06/Oct/2015:15:52:43 +0900] "GET / HTTP/1.1" 200 133 "-" "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1"
```

Pass nssh's stdin to remote command. Use `-i` option.

```
$ echo foo | nssh -i -t host-web-001 -t host-web-002 cat
foo
foo

```

Target hosts from file. Use `-f` option.

```
$ cat hosts.txt
host-web-001
host-web-002
$ nssh -f hosts.txt cat /var/log/nginx/access.log
12.34.567.890 - - [06/Oct/2015:15:52:42 +0900] "-" 400 0 "-" "-"
12.34.567.890 - - [06/Oct/2015:15:52:43 +0900] "GET / HTTP/1.1" 200 133 "-" "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1"
```
