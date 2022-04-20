build:
	pkger && go build

.PHONY: clean

clean:
	go clean

install:
	cp dp `go env GOPATH`/bin/

compress:
	upx dp

e2e: install
	(mkdir /tmp/e2e && cd /tmp/e2e && dp init)

release: clean build compress
