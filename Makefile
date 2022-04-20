build:
	pkger && go build

.PHONY: clean

clean:
	go clean

install:
	cp dp `go env GOPATH`/bin/

compress:
	upx dp

release: clean buld
