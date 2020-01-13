export GO111MODULE=on

all: bin/bento

bin/bento: cmd/bento/main.go mirait/*.go config/*.go
	go mod tidy
	go build -o bin/bento cmd/bento/main.go

vet:
	go vet ./...

errcheck:
	errcheck ./...

staticcheck:
	staticcheck -checks="all,-ST1000" ./...

clean:
	rm -rf bin/*

.PHONY: all vet errcheck staticcheck clean
