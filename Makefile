build: clean
	pkger && go build -o build/dp main.go

.PHONY: clean

clean:
	rm -rf build

linux_amd64:
	pkger && env GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o build/dp-linux_amd64 main.go
darwin_amd64:
	pkger && env GOOS=darwin GOARCH=amd64 go build -ldflags "-w" -o build/dp-darwin_amd64 main.go
windows_amd64:
	pkger && env GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o build/dp-windows_amd64.exe main.go
compress:
	upx build/dp-* 

release: clean linux_amd64 darwin_amd64 windows_amd64
