build:
	go build

build-ci: clean install-dependencies build coverage e2e compress

lint:
	golangci-lint run

ut:
	go test -v ./...

coverage:
	go test ./... -v -coverpkg=./... -coverprofile=cover.out  -p 1 ./...
	go tool cover -func=cover.out    

install-dependencies:
	go mod download

.PHONY: clean

clean:
	go clean

install:
	cp dpa `go env GOPATH`/bin/

compress:
	upx dpa

e2e: install
	(mkdir /tmp/e2e && cd /tmp/e2e && dpa init)

embed-files:
	(cd template && find init -type f -not -path "*.pyc" -not -path "*__pycache__*"  -not -path "*/venv/*")
