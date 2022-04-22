build:
	go build

lint:
	golangci-lint run

.PHONY: clean

clean:
	go clean

install:
	cp dp `go env GOPATH`/bin/

compress:
	upx dp

e2e: install
	(mkdir /tmp/e2e && cd /tmp/e2e && dp init)

embed-files:
	(cd template && find init -type f -not -path "*.pyc" -not -path "*__pycache__*"  -not -path "*/venv/*")

release: clean build compress
