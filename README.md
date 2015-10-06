# nssh

nssh is a tool to execute command in parallel on multiple hosts and returns the aggregated result.

## Installation

```
go get github.com/fujiwara/nssh
```

or

```
## at Linux
curl -slO https://github.com/fujiwara/nssh/releases/download/v0.0.1/nssh-v0.0.1-linux-amd64.zip
unzip nssh-v0.0.1-linux-amd64.zip
mv nssh-v0.0.1-linux-amd64 /usr/local/bin/nssh
```


## Useage

```
$ nssh --help
Usage of nssh:
-p=false: add hostname to line prefix
-t=[]: target hostname
-v=false: show version
```

## Sample

```
$ nssh -t host-web-001 -t host-web-002 cat /var/log/nginx/access.log
12.34.567.890 - - [06/Oct/2015:15:52:42 +0900] "-" 400 0 "-" "-"
12.34.567.890 - - [06/Oct/2015:15:52:43 +0900] "GET / HTTP/1.1" 200 133 "-" "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1"
```
