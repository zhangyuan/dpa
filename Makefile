build:
	go build

build-ci: clean install-dependencies build coverage e2e compress

lint:
	golangci-lint run

ut:
	gotest -v ./...

coverage:
	gotest ./... -v -coverpkg=./... -coverprofile=cover.out  -p 1 ./...
	go tool cover -func=cover.out    

install-dependencies:
	go install github.com/rakyll/gotest && go mod download

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
