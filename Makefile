GIT_VER := $(shell git describe --tags)

.PHONY: test packages clean

all: nssh

nssh: nssh.go
	go build -ldflags "-X main.version ${GIT_VER}"

test:
	go test

packages:
	gox -os="linux darwin" -arch="amd64" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-X main.version ${GIT_VER}"
	cd pkg && find . -name "*${GIT_VER}*" -type f -exec zip {}.zip {} \;

clean:
	rm -f pkg/* nssh
